package instagram

import (
	"fmt"
)

const (
	MessageTypeText      MessageType = "text"
	MessageTypeLike      MessageType = "like"
	MessageTypeActionLog MessageType = "action_log"
	MessageTypeLink      MessageType = "link"
	MessageTypeMedia     MessageType = "media"
	MessageTypeUndefined MessageType = "undefined"

	MediaTypeImage       MediaType = "image"
	MediaTypeVideo       MediaType = "video"
	MediaTypeVisualImage MediaType = "visual_image"
	MediaTypeVisualVideo MediaType = "visual_video"
	MediaTypeAnimated    MediaType = "animated"
	MediaTypeVoice       MediaType = "voice"
	MediaTypeUndefined   MediaType = "undefined"
)

type MessageType string

type MediaType string

type ThreadItem struct {
	ID            string
	UserID        string
	Timestamp     int64
	ClientContext string
	Type          MessageType
	Text          Text
	Media         Media
	Link          Link
}

type Text string

type Link struct {
	Url             string
	Title           string
	Summary         string
	ImagePreviewUrl string
}

type Media struct {
	ID     string
	Type   MediaType
	Width  int
	Height int
	Url    string
}

func (i ThreadItem) Validate() error {
	if i.ID == "" {
		return fmt.Errorf("ID should not be empty")
	}

	if i.UserID == "" {
		return fmt.Errorf("UserID should not be empty")
	}

	if i.Timestamp == 0 {
		return fmt.Errorf("Timestamp should not be empty")
	}

	if i.Type == "" {
		return fmt.Errorf("Type should not be empty")
	}

	return nil
}
