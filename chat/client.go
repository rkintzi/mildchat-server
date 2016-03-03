package chat

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
)

func NewClient(ws *websocket.Conn, server *server) *client {
	return &client{
		ws,
		server,
		make(chan *Message, 100),
		make(chan int),
	}
}

type client struct {
	ws       *websocket.Conn
	server   *server
	messages chan *Message
	exit     chan int
}

func (c *client) SendMessage(m *Message) {
	select {
	case c.messages <- m:
	default:
		c.Close()
	}
}

func (c *client) Receive() {
	log.Println("Receiving")
	for {
		m := &Message{}
		err := websocket.JSON.Receive(c.ws, m)
		if err == nil {
			c.server.BroadcastMessage(m)
		} else if err == io.EOF {
			c.Close()
			break
		} else {
			log.Println(err)
			c.Close()
			break
		}
	}
}

func (c *client) Close() {
	c.exit <- 1
}

func (c *client) Listen() {
	go c.Receive()

loop:
	for {
		select {
		case m := <-c.messages:
			websocket.JSON.Send(c.ws, m)
		case <-c.exit:
			c.ws.Close()
			log.Println("Closed connection")
			break loop
		}
	}
}
