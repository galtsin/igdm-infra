package channels

import (
	"fmt"
)

const (
	MessageTypeText      MessageType = "text"
	MessageTypeMedia     MessageType = "media"
	MessageTypeUndefined MessageType = "undefined"
)

type MessageType string

type MediaType string

type Message struct {
	ID             string
	AccountID      string
	ConversationID string
	Type           MessageType
	Text           string
	Media          Media
}

type Media struct {
	ID  string
	Url string
}

func (m Message) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("ID should not be empty")
	}

	if m.AccountID == "" {
		return fmt.Errorf("AccountID should not be empty")
	}

	if m.ConversationID == "" {
		return fmt.Errorf("ConversationID should not be empty")
	}

	if m.Type == "" {
		return fmt.Errorf("Type should not be empty")
	}

	return nil
}
