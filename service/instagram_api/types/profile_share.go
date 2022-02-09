package types

import (
	"fmt"

	"channels-instagram-dm/domain/model/instagram"
)

type Profile struct {
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	ProfilePicUrl string `json:"profile_pic_url"`
}

func (m Profile) ToModel() (instagram.Link, error) {
	model := instagram.Link{}
	model.Url = "https://instagram.com/" + m.Username
	model.Title = "User profile"
	model.Summary = fmt.Sprintf("%s @%s, more %s", model.Title, m.Username, model.Url)
	model.ImagePreviewUrl = m.ProfilePicUrl

	return model, nil
}
