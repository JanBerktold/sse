package main

import (
	"fmt"
	"github.com/longsleep/sse"
	"net/http"
	"time"
)

func HandleSSE(w http.ResponseWriter, r *http.Request) {
	conn, err := sse.Upgrade(w, r)

	if err != nil {
		// log error to console
		fmt.Printf("Error occured: %q", err.Error())
	}

	for {
		// Trigger event "time" with current time
		conn.WriteStringEvent("time", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
		conn.WriteStringEvent("feed", "User XY did Z")
		time.Sleep(1 * time.Second)
	}
}

func main() {

	// handle server-sent events request
	http.HandleFunc("/event", HandleSSE)

	// serve HTML page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "main.html")
	})

	http.ListenAndServe(":80", nil)
}
