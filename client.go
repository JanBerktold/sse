package sse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var (
	// ErrNilChan will be returned by Notify if it is passed a nil channel
	ErrNilChan = fmt.Errorf("nil channel given")
)

// Client is the default client used for requests.
var Client = &http.Client{}

func makeReqest(method, uri string, body io.Reader) (*http.Request, error) {
	req, err := GetReq(method, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", sseContentType)

	return req, nil
}

// Event is a go representation of an http server-sent event
type Event struct {
	URI  string
	Type string
	Data io.Reader
}

// GetReq is a function to return a single request. It will be used by notify to
// get a request and can be replaced if additional configuration is desired on
// the request. The "Accept" header will necessarily be overwritten.
var GetReq = func(method, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, uri, body)
}

// Notify takes the uri of an SSE stream and channel, and will send an Event
// down the channel when recieved, until the stream is closed. It will then
// close the stream. This is blocking, and so you will likely want to call this
// in a new goroutine (via `go Notify(..)`)
func Notify(uri string, evCh chan<- *Event) error {
	if evCh == nil {
		return ErrNilChan
	}

	req, err := makeReqest(http.MethodGet, uri, nil)
	if err != nil {
		return fmt.Errorf("error getting sse request: %v", err)
	}

	res, err := Client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing sse request for %s: %v", uri, err)
	}

	br := bufio.NewReader(res.Body)
	defer res.Body.Close()

	delim := []byte{':', ' '}

	var currEvent *Event

	for {
		bs, err := br.ReadBytes('\n')

		if err != nil && err != io.EOF {
			return err
		}

		if len(bs) < 2 {
			continue
		}

		spl := bytes.Split(bs, delim)

		if len(spl) < 2 {
			continue
		}

		currEvent = &Event{URI: uri}
		switch string(spl[0]) {
		case keyEvent:
			currEvent.Type = string(bytes.TrimSpace(spl[1]))
		case keyData:
			currEvent.Data = bytes.NewBuffer(bytes.TrimSpace(spl[1]))
			evCh <- currEvent
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}
