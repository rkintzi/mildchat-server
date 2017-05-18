package main

import (
	"encoding/json"
	"fmt"
)

type MessageType int

const (
	UnknownMessageType MessageType = iota
	ErrorMessageType
	ChatMessageType
	NickMessageType
)

func (t MessageType) String() string {
	switch t {
	case ChatMessageType:
		return "ChatMessage"
	case NickMessageType:
		return "NickMessage"
	case ErrorMessageType:
		return "ErrorMessage"
	default:
		return fmt.Sprintf("UnsupportedMessage(%d)", int(t))
	}
}

func (t MessageType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.String() + "\""), nil
}

func (t *MessageType) UnmarshalJSON(bs []byte) error {
	mt := string(bs)
	if mt == "\"ChatMessage\"" {
		*t = ChatMessageType
	} else {
		*t = UnknownMessageType
	}
	return nil
}

type Message interface {
	Type() MessageType
}

type Frame struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

const (
	ErrInvalidNick = "InvalidNick"
	ErrNoNickSet   = "NoNickSet"
)

type ErrorMessage struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

func (m *ErrorMessage) Type() MessageType {
	return ErrorMessageType
}

type ChatMessage struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (m *ChatMessage) Type() MessageType {
	return ChatMessageType
}

type NickMessage struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

func (m *NickMessage) Type() MessageType {
	return NickMessageType
}
