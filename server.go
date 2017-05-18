package main

import (
	"log"
	"net/http"
	"reflect"

	"golang.org/x/net/websocket"
)

type server struct {
	path     string
	addr     string
	clients  []*client
	reg      chan *client
	clinames map[*client]string
	nicks    map[string]bool
}

func NewServer(path, addr string) *server {
	return &server{
		path:     path,
		addr:     addr,
		clients:  make([]*client, 0, 10),
		reg:      make(chan *client),
		clinames: make(map[*client]string),
		nicks:    make(map[string]bool),
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
			message := value.Interface().(Message)
			s.handleMessage(message, s.clients[chosen-1])
		}
	}
}
func (s *server) handleMessage(m Message, sender *client) {
	switch msg := m.(type) {
	case *ChatMessage:
		msg.Author = s.clinames[sender]
		if msg.Author == "" {
			log.Printf("Ignoring client %v: no nick name", sender.Id())
			errMsg := ErrorMessage{
				ErrorCode: ErrNoNickSet,
				Message:   "Set your nick name first",
			}
			sender.Send(&errMsg)
			break
		}
		log.Printf("Broadcast chat message from: %v\n", sender.Id())
		s.broadcast(msg, sender)
	case *NickMessage:
		if msg.NewName == "" || s.nicks[msg.NewName] {
			log.Printf("Ignoring client %v: invalid nick or nick already taken: %s", sender.Id(), msg.NewName)
			errMsg := ErrorMessage{
				ErrorCode: ErrInvalidNick,
				Message:   "Invalid nick name or nick name already taken",
			}
			sender.Send(&errMsg)
			break
		}
		msg.OldName = s.clinames[sender]
		if msg.OldName != "" {
			s.nicks[msg.OldName] = false
		}
		s.nicks[msg.NewName] = true
		s.clinames[sender] = msg.NewName
		log.Printf("Broadcast nick message from: %v\n", sender.Id())
		s.broadcast(msg, sender)
	}
}
func (s *server) broadcast(m Message, sender *client) {
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
