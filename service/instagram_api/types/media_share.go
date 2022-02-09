package types

import "channels-instagram-dm/domain/model/instagram"

type MediaShare struct {
	Media
}

func (m MediaShare) ToModel() (instagram.Media, error) {
	return m.Media.ToModel()
}
