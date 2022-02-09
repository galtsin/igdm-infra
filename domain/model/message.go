package model

import (
	"time"
)

const (
	MessageTypeText             MessageType = "text"
	MessageTypeLike             MessageType = "like"
	MessageTypeLink             MessageType = "link"
	MessageTypeActionLog        MessageType = "action_log"
	MessageTypeMediaImage       MessageType = "media_image"
	MessageTypeMediaVideo       MessageType = "media_video"
	MessageTypeMediaVisualImage MessageType = "media_visual_image"
	MessageTypeMediaVisualVideo MessageType = "media_visual_video"
	MessageTypeMediaAnimated    MessageType = "media_animated"
	MessageTypeMediaVoice       MessageType = "media_voice"
	MessageTypeUndefined        MessageType = "undefined"
)

const (
	MessageSourceInstagram MessageSource = "instagram"
	MessageSourceChannels  MessageSource = "channels"
)

const (
	MessageDeliveryStatusNone MessageDeliveryStatus = iota
	MessageDeliveryStatusWaiting
	MessageDeliveryStatusSuccess
	MessageDeliveryStatusFailed
)

type MessageType string

type MessageSource string

type MessageDeliveryStatus uint

type Message struct {
	ID             string
	AccountID      string
	ConversationID string
	Source         MessageSource
	Type           MessageType
	Payload        interface{}
	Attributes     MessageAttributes
	Delivered      MessageDelivered
	CreatedAt      time.Time
}

type MessageAttributes struct {
	InstagramAttributes
	ChannelsAttributes
}

type InstagramAttributes struct {
	ID        string
	UserID    string
	Timestamp int64
}

type ChannelsAttributes struct {
	ID string
}

type MessageDelivered struct {
	Status    MessageDeliveryStatus
	AttemptAt time.Time
}

type MessageUndefined struct {
	Text string
}

type MessageLike struct {
	Like string
}

type MessageText struct {
	Text string
}

type MessageLink struct {
	Url             string
	Title           string
	Summary         string
	ImagePreviewUrl string
}

type MessageActionLog struct {
	Text string
}

type MessageMedia struct {
	ID  string
	Url string
}

type MessageMediaImage struct {
	MessageMedia
	Width  int
	Height int
}

type MessageMediaVideo struct {
	MessageMedia
	Width  int
	Height int
}

type MessageMediaVisualImage struct {
	MessageMedia
	Width  int
	Height int
}

type MessageMediaVisualVideo struct {
	MessageMedia
	Width  int
	Height int
}

type MessageMediaAnimated struct {
	MessageMedia
	Width  int
	Height int
}

type MessageMediaVoice struct {
	MessageMedia
}

type MessagesBatch map[string][]Message

func NewMessage(accountID string, conversationID string, source MessageSource) Message {
	return Message{
		AccountID:      accountID,
		ConversationID: conversationID,
		Source:         source,
		Attributes:     MessageAttributes{},
		CreatedAt:      time.Now(),
	}
}

func (m *Message) SetInstagramAttributes(attributes InstagramAttributes) {
	m.Attributes.InstagramAttributes = attributes
}

func (m *Message) SetChannelsAttributes(attributes ChannelsAttributes) {
	m.Attributes.ChannelsAttributes = attributes
}

func (m *Message) DeliveredWaiting() {
	m.Delivered.Status = MessageDeliveryStatusWaiting
	m.Delivered.AttemptAt = time.Now()
}

func (m *Message) DeliveredSuccess() {
	m.Delivered.Status = MessageDeliveryStatusSuccess
	m.Delivered.AttemptAt = time.Now()
}

func (m *Message) DeliveredFail() {
	m.Delivered.Status = MessageDeliveryStatusFailed
	m.Delivered.AttemptAt = time.Now()
}

func (m *Message) SetPayload(payload interface{}) {
	switch payload.(type) {
	case MessageText:
		m.Type = MessageTypeText
	case MessageLike:
		m.Type = MessageTypeLike
	case MessageLink:
		m.Type = MessageTypeLink
	case MessageActionLog:
		m.Type = MessageTypeActionLog
	case MessageMediaImage:
		m.Type = MessageTypeMediaImage
	case MessageMediaVideo:
		m.Type = MessageTypeMediaVideo
	case MessageMediaVisualImage:
		m.Type = MessageTypeMediaVisualImage
	case MessageMediaVisualVideo:
		m.Type = MessageTypeMediaVisualVideo
	case MessageMediaAnimated:
		m.Type = MessageTypeMediaAnimated
	case MessageMediaVoice:
		m.Type = MessageTypeMediaVoice
	case MessageUndefined:
		m.Type = MessageTypeUndefined
	default:
		m.Type = MessageTypeUndefined
	}

	m.Payload = payload
}
