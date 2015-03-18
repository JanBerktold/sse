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

func (c *Conn) Write(msg []byte) error {
	if !c.isOpen {
		return ErrConnectionClosed
	} else {
		c.messages <- msg
		return nil
	}
}

func (c *Conn) WriteString(msg string) error {
	return c.Write([]byte(msg))
}

func (c *Conn) WriteJson(value interface{}) error {
	if by, err := json.Marshal(value); err == nil {
		return c.Write(by)
	} else {
		return err
	}
}

func (c *Conn) WriteXml(value interface{}) error {
	if by, err := xml.Marshal(value); err == nil {
		return c.Write(by)
	} else {
		return err
	}
}
