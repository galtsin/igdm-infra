package types

import (
	"fmt"

	"channels-instagram-dm/domain/model/instagram"
)

type ReelShare struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Media struct {
		Media
		PK      interface{} `json:"pk"`
		User    User        `json:"user"`
		Caption struct {
			PK        interface{} `json:"pk"`
			UserID    interface{} `json:"user_id"`
			Text      string      `json:"text"`
			Type      int         `json:"type"`
			User      User        `json:"user"`
			IsCovered bool        `json:"is_covered"`
			MediaID   interface{} `json:"media_id"`
		} `json:"caption"`
	} `json:"media"`
}

func (m ReelShare) ToModel() (instagram.Link, error) {
	model := instagram.Link{}
	model.Url = fmt.Sprintf("https://instagram.com/stories/%s/%s", m.Media.User.Username, ValueToString(m.Media.PK))

	switch m.Type {
	case "reply":
		model.Title = "Story comment"
		model.Summary = fmt.Sprintf("%s @%s: %s, more %s", model.Title, m.Media.User.Username, m.Text, model.Url)
	case "mention":
		model.Title = "Story mention"
		model.Summary = fmt.Sprintf("%s @%s: %s, more %s", model.Title, m.Media.User.Username, m.Media.Caption.Text, model.Url)
	default:
		model.Title = fmt.Sprintf("Unsupported reel_share type %s", m.Type)
		model.Summary = fmt.Sprintf("%s, more %s", model.Title, model.Url)
	}

	return model, nil
}
