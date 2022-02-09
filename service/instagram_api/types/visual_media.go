package types

import (
	"channels-instagram-dm/domain/model/instagram"
)

type VisualMedia struct {
	Media           Media `json:"media"`
	UrlExpireAtSecs int64 `json:"url_expire_at_secs"`
}

func (m VisualMedia) ToModel() (instagram.Media, error) {
	model := instagram.Media{
		ID:   ValueToString(m.Media.ID),
		Type: instagram.MediaTypeUndefined,
	}

	switch m.Media.MediaType {
	// photo
	case 1:
		best := m.Media.Images.GetBestCandidate()

		model.Width = best.Width
		model.Height = best.Height
		model.Url = best.Url
		model.Type = instagram.MediaTypeVisualImage
	// video
	case 2:
		best := m.Media.Videos.GetBestCandidate()

		model.Width = best.Width
		model.Height = best.Height
		model.Url = best.Url
		model.Type = instagram.MediaTypeVisualVideo
	// carousel
	// case 8:
	default:
		model.Type = instagram.MediaTypeUndefined
	}

	return model, nil
}
