package types

import "channels-instagram-dm/domain/model/instagram"

type StoryShare struct {
	Media `json:"media"`
}

func (m StoryShare) ToModel() (instagram.Media, error) {
	return m.Media.ToModel()
}
