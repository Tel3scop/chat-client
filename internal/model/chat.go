package model

import "github.com/golang/protobuf/ptypes/timestamp"

// ChatMessage структура сообщений
type ChatMessage struct {
	ChatID    int64
	From      string
	Text      string
	Timestamp *timestamp.Timestamp
}

// Chat структура чата
type Chat struct {
	ID      int64
	Name    string
	Members []string
}
