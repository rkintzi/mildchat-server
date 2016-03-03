package chat

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

func NewServer(path, addr string) *server {
	return &server{
		path,
		addr,
		make([]*client, 0),
	}
}

type server struct {
	path    string
	addr    string
	clients []*client
}

func (s *server) BroadcastMessage(m *Message) {
	log.Println("Sending message")
	for _, c := range s.clients {
		c.SendMessage(m)
	}
}

func (s *server) connectHandler(ws *websocket.Conn) {
	log.Println("Client connected")
	c := NewClient(ws, s)
	s.clients = append(s.clients, c)

	c.Listen()

	log.Println("Client disconnected")
	for i, cl := range s.clients {
		if cl == c {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}
}

func (s *server) Listen() {
	log.Println("Starting server")

	http.Handle(s.path, websocket.Handler(s.connectHandler))

	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		log.Panic("Server error: " + err.Error())
	}
}
