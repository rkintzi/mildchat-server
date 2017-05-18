package main

import "log"

func main() {
	server := NewServer("/chat", ":8080")
	log.Println("Listening...")
	if err := server.Listen(); err != nil {
		log.Fatalf("Listen error: %v\n", err)
	}
}
