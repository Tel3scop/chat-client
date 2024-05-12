package model

import "github.com/golang/protobuf/ptypes/timestamp"

type ChatMessage struct {
	ChatID    int64
	From      string
	Text      string
	Timestamp *timestamp.Timestamp
}

type Chat struct {
	ID      int64
	Name    string
	Members []string
}
