package main

import (
	"github.com/lwojciechowski/mildchat-server/chat"
)

func main() {
	server := chat.NewServer("/chat", ":8080")
	server.Listen()
}
