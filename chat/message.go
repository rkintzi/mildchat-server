package chat

type Message struct {
	Author  string `json:"author"`
	Message string `json:"message"`
}
