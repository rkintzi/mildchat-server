package chat

import (
	"code.google.com/p/go.net/websocket"
)

func NewClient(ws *websocket.Conn, server *server) *client {
	return &client{
		ws,
		server,
		make(chan *Message),
	}
}

type client struct {
	ws       *websocket.Conn
	server   *server
	messages chan *Message
}

func (self *client) SendMessage(message *Message) {
	self.messages <- message
}

func (self *client) Listen() {
	for {
		select {}
	}
}
