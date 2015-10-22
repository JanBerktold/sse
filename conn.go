package sse

import (
	"encoding/json"
	"encoding/xml"
	"time"
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
	return c.WriteEvent("", msg)
}

// Sends a byte-slice to the connected client and triggers the specified event with the data. Returns an error
// if the connection is already closed.
func (c *Conn) WriteEvent(typ string, msg []byte) error {
   for {
      select {
      case c.messages <- message{
         message: msg,
         typ:     typ,
      }:
         return nil
      default:
         if !c.isOpen {
            return ErrConnectionClosed
         }
         time.Sleep(10*time.Microsecond) // give channel time to unblock
      }
   }
}

// Sends a string to the connected client. Returns an error
// if the connection is already closed.
func (c *Conn) WriteString(msg string) error {
	return c.WriteEvent("", []byte(msg))
}

// Sends a string to the connected client, targeting the specified event. Returns an error
// if the connection is already closed.
func (c *Conn) WriteStringEvent(typ, msg string) error {
	return c.WriteEvent(typ, []byte(msg))
}

// Sends a json-encoded struct to the connected client. Returns an error
// if the connection is already closed or if the encoding failed.
func (c *Conn) WriteJson(value interface{}) error {
	return c.WriteJsonEvent("", value)
}

// Sends a json-encoded struct to the connected client, targeting the specified event. Returns an error
// if the connection is already closed or if the encoding failed.
func (c *Conn) WriteJsonEvent(typ string, value interface{}) error {
	if by, err := json.Marshal(value); err == nil {
		return c.WriteEvent(typ, by)
	} else {
		return err
	}
}

// Sends a xml-encoded struct to the connected client. Returns an error
// if the connection is already closed or if the encoding failed.
func (c *Conn) WriteXml(value interface{}) error {
	return c.WriteXmlEvent("", value)
}

// Sends a xml-encoded struct to the connected client, targeting the specified event. Returns an error
// if the connection is already closed or if the encoding failed.
func (c *Conn) WriteXmlEvent(typ string, value interface{}) error {
	if by, err := xml.Marshal(value); err == nil {
		return c.WriteEvent(typ, by)
	} else {
		return err
	}
}

// Returns whether the connection is still opened.
func (c *Conn) IsOpen() bool {
	return c.isOpen
}

// Forces the connection to close. The Conn object should not be used anymore.
func (c *Conn) Close() {
	c.shutdown <- true
}
