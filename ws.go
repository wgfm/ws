package ws

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const wsUUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

var (
	ErrNotUpgradable = errors.New("connection not upgradable")
)

func Upgrade(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	if !strings.EqualFold(r.Header.Get("Connection"), "upgrade") {
		return nil, fmt.Errorf("%w: connection header != upgrade", ErrNotUpgradable)
	}

	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return nil, fmt.Errorf("%w: upgrade header != websocket", ErrNotUpgradable)
	}

	h, ok := w.(http.Hijacker)
	if !ok {
		return nil, fmt.Errorf("%w: connection not hijackable", ErrNotUpgradable)
	}

	wsKey := r.Header.Get("Sec-WebSocket-Key")
	w.Header().Set("Sec-WebSocket-Accept", wsAccept(wsKey))
	w.WriteHeader(http.StatusSwitchingProtocols)
	w.Header().Set("Upgrade", "websocket")

	conn, bufrw, err := h.Hijack()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotUpgradable, err)
	}

	bufrw.Flush()
	return conn, nil
}

func wsAccept(wsKey string) string {
	s := sha1.Sum([]byte(wsKey + wsUUID))
	return base64.StdEncoding.EncodeToString(s[:])
}
