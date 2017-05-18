package main

import (
	"log"
	"net/http"
	"reflect"

	"golang.org/x/net/websocket"
)

type server struct {
	path    string
	addr    string
	clients []*client
	reg     chan *client
}

func NewServer(path, addr string) *server {
	return &server{
		path:    path,
		addr:    addr,
		clients: make([]*client, 0, 10),
		reg:     make(chan *client),
	}
}

func (s *server) handleClients() {
	cases := make([]reflect.SelectCase, 1)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(s.reg)}
	nCli := 0
	for !cases[0].Chan.IsNil() || nCli > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			if chosen == 0 {
				// Listen channel was closed; signal clients server is about to be closed
				// @@rk: should we add timeout to not wait for clients forever?
				s.stop()
			} else {
				log.Printf("Client disconnected: %v\n", s.clients[chosen-1].Id())
				s.clients[chosen-1] = nil
			}
			// TODO clean s.clients and cases lists every n-time we are here
			continue
		}

		if chosen == 0 {
			// New client connected
			cli := value.Interface().(*client)
			s.clients = append(s.clients, cli)
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cli.messages)})
			nCli++
		} else {
			if cases[0].Chan.IsNil() {
				log.Println("Message from: %v (in shutdown mode - ignored)\n", s.clients[chosen-1].Id())
				continue
			}
			message := value.Interface().(*ChatMessage)
			s.broadcast(message, s.clients[chosen-1])
		}
	}
}

func (s *server) broadcast(m *ChatMessage, sender *client) {
	log.Printf("Broadcast message from: %v\n", sender.Id())
	for _, c := range s.clients {
		if c == nil {
			continue
		}
		err := c.Send(m)
		if err != nil {
			log.Printf("Can't send to client: %v: %v", c.Id(), err)
		}
	}
}

func (s *server) stop() {
	for _, c := range s.clients {
		c.Close()
	}
}

func (s *server) connectHandler(c *websocket.Conn) {
	cli := NewClient(c)
	log.Printf("New client connected: %v\n", cli.Id())
	s.reg <- cli
	cli.Receive()
}

func (s *server) Listen() error {
	go s.handleClients()
	http.Handle(s.path, websocket.Handler(s.connectHandler))
	return http.ListenAndServe(s.addr, nil)
}
