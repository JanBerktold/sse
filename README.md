# Server-sent events for Go
[![GoDoc](https://godoc.org/github.com/longsleep/sse?status.svg)](https://godoc.org/github.com/longsleep/sse)

This is a lightweight SEE library for Golang which is designed to play nicely along different packages and provide a convient usage. Compatible with every Go version since 1.1.

## Examples

*Note:* Also look into the examples folder.

```go
import (
	"net/http"
	"github.com/longsleep/sse"
)

type Person struct {
	Name string
	Age int
}

func main() {

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		// get a SSE connection from the HTTP request
		// in a real world situation, you should look for the error (second return value)
		conn, _ := sse.Upgrade(w, r)

		// writes the struct as JSON
		conn.WriteJson(&Person{
			Name: "Jan",
			Age: 17,
		})

		// writes the struct as XML
		conn.WriteXml(&Person{
			Name: "Ernst",
			Age: 23,
		})

		// write a plain string
		conn.WriteString("Hello how are you?")

		// trigger the event "feed" with "This is a message" as payload
		// [extended example](https://github.com/longsleep/sse/tree/master/examples/events)
		conn.WriteStringEvent("feed", "This is a message")

		for {
			// keep this goroutine alive to keep the connection

			// get a message from some channel
			// blocks until it recieves a messages and then instantly sends it to the client
			msg := <-someChannel
			conn.WriteString(msg)
		}
	})

	http.ListenAndService(":80", nil)
}

```

### Usage with a Upgrader instance

Using a Upgrader instance allows you to specify a RetryTime interval at which the client will attempt to reconnect to the EventSource.

```go
import (
	"net/http"
	"github.com/longsleep/sse"
)

type Person struct {
	Name string
	Age int
}

func main() {

	upgrader := sse.Upgrader{
		RetryTime: 5 * time.Second,
	}

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		// get a SSE connection from the HTTP request
		// in a real world situation, you should look for the error (second return value)
		conn, _ := upgrader.Upgrade(w, r)

		// writes the struct as JSON
		conn.WriteJson(&Person{
			Name: "Jan",
			Age: 17,
		})

		for {
			// get a message from some channel
			// blocks until it recieves a messages and then instantly sends it to the client
			msg := <-someChannel
			conn.WriteString(msg)
		}
	})

	http.ListenAndService(":80", nil)
}

```
