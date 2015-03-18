package sse

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrStreamingNotSupported = errors.New("Streaming unsupported!")
	ErrConnectionClosed      = errors.New("Connection already closed")
)

func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return nil, ErrStreamingNotSupported
	}

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	conn := &Conn{
		messages: make(chan []byte),
		shutdown: make(chan bool),
		isOpen:   true,
	}

	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		for {
			select {
			case msg := <-conn.messages:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				f.Flush()
			case <-conn.shutdown:
				conn.isOpen = false
				return
			case <-notify:
				conn.isOpen = false
				return
			}
		}
	}()

	return conn, nil
}
