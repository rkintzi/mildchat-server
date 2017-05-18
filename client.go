package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

// NewClient is client constructor
func NewClient(ws *websocket.Conn) *client {
	return &client{
		ws,
		make(chan Message, 100),
	}
}

type client struct {
	con      *websocket.Conn
	messages chan Message
}

func (c *client) Id() string {
	return fmt.Sprintf("%v", c.con.Request().RemoteAddr)
}

func (c *client) Send(m Message) error {
	var err error
	f := Frame{Type: m.Type()}
	f.Data, err = json.Marshal(m)
	if err != nil {
		return err
	}
	return websocket.JSON.Send(c.con, f)
}

// Receive reads websocket in the loop.
func (c *client) Receive() {
	log.Printf("Start receiving from: %v\n", c.Id())
loop:
	for {
		f := &Frame{}
		err := websocket.JSON.Receive(c.con, f)
		if err == io.EOF {
			close(c.messages)
			break
		} else if err != nil {
			log.Printf("Error from client %v: %v", c.Id(), err)
			close(c.messages)
			break
		}
		var msg Message
		switch f.Type {
		case ChatMessageType:
			msg = &ChatMessage{}
		case NickMessageType:
			msg = &NickMessage{}
		default:
			log.Printf("Unsupported message type from client %v: %s", c.Id(), f.Type)
			close(c.messages)
			break loop
		}
		err = json.Unmarshal(f.Data, &msg)
		if err != nil {
			log.Printf("Can not unparse data %v: %v: %v", c.Id(), err, f.Data)
			close(c.messages)
			break
		}
		c.messages <- msg
	}
}

func (c *client) Close() {
	c.con.Close()
}
