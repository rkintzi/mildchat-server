package chat

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

func NewServer(path string) *server {
	return &server{
		path,
		make([]*Client, 0),
		make([]*Message, 0),
		make(chan *Client),
		make(chan *Client),
		make(chan *Message),
	}
}

type server struct {
	path         string
	clients      []*Client
	messages     []*Message
	addClient    chan *Client
	removeClient chan *Client
	sendMessage  chan *Message
}

func (self *server) AddClient(client *Client) {
	log.Println("Adding client")
	self.addClient <- client
}

func (self *server) RemoveClient(client *Client) {
	log.Println("Removing client")
	self.removeClient <- client
}

func (self *server) SendMessage(message *Message) {
	log.Println("Sending message")
	self.sendMessage <- message
}

func (self *server) connectHandler(ws *websocket.Conn) {
	client := NewClient(ws, self)
	self.AddClient(client)
	client.Listen()
	defer ws.Close()
}

func (self *server) Listen() {
	log.Println("Starting server")
}
