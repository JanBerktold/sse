package sse

import (
	"encoding/json"
	"encoding/xml"
)

type Conn struct {
	messages chan []byte
	shutdown chan bool
	isOpen   bool
}

// Sends a byte-slice to the connected client. Returns an error
// if the connection is already closed.
func (c *Conn) Write(msg []byte) error {
	if !c.isOpen {
		return ErrConnectionClosed
	} else {
		c.messages <- msg
		return nil
	}
}

// Sends a string to the connected client. Returns an error
// if the connection is already closed.
func (c *Conn) WriteString(msg string) error {
	return c.Write([]byte(msg))
}

// Sends a json-encoded struct to the connected client. Returns an error
// if the connection is already closed or if the encoding failed.
func (c *Conn) WriteJson(value interface{}) error {
	if by, err := json.Marshal(value); err == nil {
		return c.Write(by)
	} else {
		return err
	}
}

// Sends a xml-encoded struct to the connected client. Returns an error
// if the connection is already closed or if the encoding failed.
func (c *Conn) WriteXml(value interface{}) error {
	if by, err := xml.Marshal(value); err == nil {
		return c.Write(by)
	} else {
		return err
	}
}

// Returns whether the connection is still opened.
func (c *Conn) IsOpen() bool {
	return c.IsOpen()
}

// Forces the connection to close. The Conn object should not be used anymore.
func (c *Conn) Close() {
	c.shutdown <- true
}
