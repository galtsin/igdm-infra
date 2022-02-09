package mq

import (
	"encoding/json"
	"time"

	"channels-instagram-dm/domain/model/channels"
	"github.com/google/uuid"
)

const (
	InboundDirection  Direction = "inbound"
	OutboundDirection Direction = "outbound"
)

const (
	ChannelsSubject = "inbound-messages"
)

type Direction string

type Packet struct {
	Uuid        string    `json:"uuid"`
	App         string    `json:"app"`
	Integration string    `json:"integration"`
	Data        Payload   `json:"data"`
	Direction   Direction `json:"direction"`
	Delivered   bool      `json:"delivered"`
	CreatedAt   time.Time `json:"created_at"`
	Error       string    `json:"error"`
}

type Payload struct {
	Message      Message      `json:"message"`
	Conversation Conversation `json:"conversation"`
	Sender       Sender       `json:"sender"`
	Timestamp    int64        `json:"timestamp"`
}

type Message struct {
	ID    string               `json:"id"`
	Type  channels.MessageType `json:"type"`
	Text  string               `json:"text"`
	Media Media                `json:"media"`
}

type Media struct {
	ID string `json:"id"`
	// Type string `json:"type"`
	Url string `json:"url"`
}

type Conversation struct {
	ID string `json:"id"`
}

type Sender struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func Marshal(appName string, dir Direction, integration string, data Payload) ([]byte, error) {
	packet := Packet{
		Uuid:        uuid.New().String(),
		App:         appName,
		Direction:   dir,
		Integration: integration,
		Data:        data,
		Delivered:   false,
		CreatedAt:   time.Now(),
	}

	payload, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func Unmarshal(data []byte) (Packet, error) {
	packet := Packet{}
	err := json.Unmarshal(data, &packet)
	return packet, err
}
