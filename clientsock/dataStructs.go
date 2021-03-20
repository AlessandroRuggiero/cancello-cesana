package clientsock

import "github.com/AlessandroRuggiero/go-socketcast"

type Message struct {
	Type int    `json:"type"`
	Body string `josn:"body"`
}

type AuthRequest struct {
	Type  int    `json:"type"`
	Token string `josn:"token"`
}

type CallbackData struct {
	Client *socketcast.Client
}
