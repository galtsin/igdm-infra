package types

import "channels-instagram-dm/domain/model/instagram"

type Link struct {
	Text        string `json:"text"`
	LinkContext struct {
		LinkUrl      string `json:"link_url"`
		LinkTitle    string `json:"link_title"`
		LinkSummary  string `json:"link_summary"`
		LinkImageUrl string `json:"link_image_url"`
	} `json:"link"`
}

func (m Link) ToModel() (instagram.Link, error) {
	model := instagram.Link{}

	model.Url = m.Text
	model.Title = m.LinkContext.LinkTitle
	model.Summary = m.LinkContext.LinkSummary
	model.ImagePreviewUrl = m.LinkContext.LinkImageUrl

	return model, nil
}
