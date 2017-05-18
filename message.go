package main

type MessageType int

const (
	ChatMessageType MessageType = iota
)

type Message interface {
	Type() MessageType
}

type ChatMessage struct {
	Author  string `json:"author"`
	Message string `json:"message"`
}

func (m *ChatMessage) Type() MessageType {
	return ChatMessageType
}
