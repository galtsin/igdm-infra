package types

import (
	"strconv"

	"channels-instagram-dm/domain/model/instagram"
)

type AnimatedMedia struct {
	ID     interface{} `json:"id"`
	Images struct {
		FixedHeight struct {
			Width    string `json:"width"`
			Height   string `json:"height"`
			Mp4      string `json:"mp4"`
			Mp4Size  string `json:"mp4_size"`
			Size     string `json:"size"`
			Url      string `json:"url"`
			Webp     string `json:"webp"`
			WebpSize string `json:"webp_size"`
		} `json:"fixed_height"`
	} `json:"images"`
}

func (m AnimatedMedia) ToModel() (instagram.Media, error) {
	model := instagram.Media{
		ID:   ValueToString(m.ID),
		Type: instagram.MediaTypeAnimated,
	}

	if m.Images.FixedHeight.Url != "" {
		w, err := strconv.ParseInt(m.Images.FixedHeight.Width, 10, 64)
		if err != nil {
			return instagram.Media{}, err
		}

		h, err := strconv.ParseInt(m.Images.FixedHeight.Height, 10, 64)
		if err != nil {
			return instagram.Media{}, err
		}

		model.Width = int(w)
		model.Height = int(h)
		model.Url = m.Images.FixedHeight.Url
	}

	return model, nil
}
