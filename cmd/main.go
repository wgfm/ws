package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/wgfm/ws"
)

func main() {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := ws.Upgrade(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %v", err)
		}

		// We're now TCPing baby
		go func() {
			io.Copy(os.Stdout, conn)
		}()

		conn.Write([]byte("Hello"))
	})

	http.ListenAndServe(":8080", nil)
}
