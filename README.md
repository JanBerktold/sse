# Server-sent events for Go
[![Build Status](https://travis-ci.org/JanBerktold/sse.svg)](https://travis-ci.org/JanBerktold/sse) [![GoDoc](https://godoc.org/github.com/JanBerktold/sse?status.svg)](https://godoc.org/github.com/JanBerktold/sse)

This is a lightweight SEE library for Golang which is designed to play nicely along different packages and provide a convient usage. Compatible with every Go version since 1.1.

## Examples

*Note:* Also look into the examples folder.

```go
import (
	"net/http"
	"github.com/janberktold/sse"
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
		// [extended example](https://github.com/JanBerktold/sse/tree/master/examples/events)
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
	"github.com/janberktold/sse"
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


## The MIT License (MIT)

Copyright (c) 2015 Jan Berktold

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

