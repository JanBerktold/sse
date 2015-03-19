package sse

import (
	"encoding/json"
	"encoding/xml"
)

type Conn struct {
	messages chan message
	shutdown chan bool
	isOpen   bool
}

type message struct {
	message []byte
	typ     string
}

// Sends a byte-slice to the connected client. Returns an error
// if the connection is already closed.
func (c *Conn) Write(msg []byte) error {
	return c.WriteEvent(msg, "")
}

func (c *Conn) WriteEvent(msg []byte, typ string) error {
	if !c.isOpen {
		return ErrConnectionClosed
	} else {
		c.messages <- message{
			message: msg,
			typ:     typ,
		}
		return nil
	}
}

// Sends a string to the connected client. Returns an error
// if the connection is already closed.
func (c *Conn) WriteString(msg string) error {
	return c.WriteEvent([]byte(msg), "")
}

func (c *Conn) WriteStringEvent(msg, typ string) error {
	return c.WriteEvent([]byte(msg), typ)
}

// Sends a json-encoded struct to the connected client. Returns an error
// if the connection is already closed or if the encoding failed.
func (c *Conn) WriteJson(value interface{}) error {
	return c.WriteJsonEvent(value, "")
}

func (c *Conn) WriteJsonEvent(value interface{}, typ string) error {
	if by, err := json.Marshal(value); err == nil {
		return c.WriteEvent(by, typ)
	} else {
		return err
	}
}

// Sends a xml-encoded struct to the connected client. Returns an error
// if the connection is already closed or if the encoding failed.
func (c *Conn) WriteXml(value interface{}) error {
	return c.WriteXmlEvent(value, "")
}

func (c *Conn) WriteXmlEvent(value interface{}, typ string) error {
	if by, err := xml.Marshal(value); err == nil {
		return c.WriteEvent(by, typ)
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
