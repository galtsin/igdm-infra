package types

import (
	"channels-instagram-dm/domain/model/instagram"
)

type VoiceMedia struct {
	Media struct {
		ID        interface{} `json:"id"`
		MediaType int         `json:"media_type"`
		Audio     struct {
			Src      string `json:"audio_src"`
			Duration int    `json:"duration"`
		} `json:"audio"`
	} `json:"media"`
}

func (m VoiceMedia) ToModel() (instagram.Media, error) {
	return instagram.Media{
		ID:   ValueToString(m.Media.ID),
		Type: instagram.MediaTypeVoice,
		Url:  m.Media.Audio.Src,
	}, nil
}
