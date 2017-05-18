package main

func main() {
	server := NewServer("/chat", ":8080")
	server.Listen()
}
