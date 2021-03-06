package sse

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrStreamingNotSupported = errors.New("Streaming unsupported!")
	ErrConnectionClosed      = errors.New("Connection already closed")

	globalUpgrader = Upgrader{}
)

type Upgrader struct {
	// time between two connects from a client
	RetryTime time.Duration
}

// Takes over a HTTP-connection and returns a SSE-Connection, which can be used
// to send events. Returns an error, if the connection does not support streaming.
// Please note, that in this case the client will also be notified and the
// HTTP-connection should therefore not be used anymore.
func (up Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {

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
		messages: make(chan message),
		shutdown: make(chan bool),
		isOpen:   true,
	}

	notify := w.(http.CloseNotifier).CloseNotify()

	// tell client about retry time
	if up.RetryTime > 0 {
		fmt.Fprintf(w, "retry: %s\n", up.RetryTime)
	}

	go func() {
		for {
			select {
			case msg := <-conn.messages:
				if len(msg.id) > 0 {
					fmt.Fprintf(w, "id: %s\n", msg.id)
				}
				if len(msg.typ) > 0 {
					fmt.Fprintf(w, "event: %s\n", msg.typ)
				}
				fmt.Fprintf(w, "data: %s\n\n", msg.message)
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

// Global Upgrade for method for usage without a Upgrader instance.
// Refer to Upgrader.Upgrade for complete documentation.
func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	return globalUpgrader.Upgrade(w, r)
}
