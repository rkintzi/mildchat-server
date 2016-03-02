package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
)

func echoHandler(ws *websocket.Conn) {
	log.Println("Connected")
	io.Copy(ws, ws)
}

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
