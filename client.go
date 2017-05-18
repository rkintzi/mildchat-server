package main

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

// NewClient is client constructor
func NewClient(ws *websocket.Conn) *client {
	return &client{
		ws,
		make(chan *ChatMessage, 100),
	}
}

type client struct {
	con      *websocket.Conn
	messages chan *ChatMessage
}

func (c *client) Id() string {
	return fmt.Sprintf("%v", c.con.Request().RemoteAddr)
}

func (c *client) Send(m *ChatMessage) error {
	return websocket.JSON.Send(c.con, m)
}

// Receive reads websocket in the loop.
func (c *client) Receive() {
	log.Printf("Start receiving from: %v\n", c.Id())
	for {
		m := &ChatMessage{}
		err := websocket.JSON.Receive(c.con, m)
		if err == io.EOF {
			close(c.messages)
			break
		} else if err != nil {
			log.Printf("Error from client %v: %v", c.Id(), err)
			close(c.messages)
			break
		}
		c.messages <- m
	}
}

func (c *client) Close() {
	c.con.Close()
}
