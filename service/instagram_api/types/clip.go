package types

import (
	"channels-instagram-dm/domain/model/instagram"
)

type Clip struct {
	Media `json:"clip"`
}

func (m Clip) ToModel() (instagram.Media, error) {
	return m.Media.ToModel()
}
