package types

import (
	"channels-instagram-dm/domain/model/instagram"
)

const (
	MediaTypePhoto       = 1
	MediaTypeVideo       = 2
	MediaTypeCarousel    = 8
	MediaTypeDirectAudio = 11
)

type Media struct {
	ID                   interface{} `json:"id"`
	Images               Images      `json:"image_versions2"`
	Videos               Videos      `json:"video_versions,omitempty"`
	OriginalWidth        int         `json:"original_width"`
	OriginalHeight       int         `json:"original_height"`
	MediaType            int         `json:"media_type"`
	MediaID              interface{} `json:"media_id"`
	PlaybackDurationSecs int         `json:"playback_duration_secs"`
	URLExpireAtSecs      int         `json:"url_expire_at_secs"`
	OrganicTrackingToken string      `json:"organic_tracking_token"`
}

type Images struct {
	Versions []Candidate `json:"candidates"`
}

type Videos []struct {
	Candidate
	ID interface{} `json:"id"`
}

type Candidate struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Url    string `json:"url"`
}

func (v Videos) GetBestCandidate() Candidate {
	best := ""
	var mh, mw int
	for _, c := range v {
		if mw < c.Width && c.Height > mh && c.Url != "" {
			mw = c.Width
			mh = c.Height
			best = c.Url
		}
	}

	return Candidate{
		Width:  mw,
		Height: mh,
		Url:    best,
	}
}

func (i Images) GetBestCandidate() Candidate {
	best := ""
	var mh, mw int
	for _, v := range i.Versions {
		if v.Width > mw || v.Height > mh {
			best = v.Url
			mh, mw = v.Height, v.Width
		}
	}

	return Candidate{
		Width:  mw,
		Height: mh,
		Url:    best,
	}
}

func (m Media) ToModel() (instagram.Media, error) {
	model := instagram.Media{
		ID:   ValueToString(m.ID),
		Type: instagram.MediaTypeUndefined,
	}

	switch m.MediaType {
	// photo
	case MediaTypePhoto:
		best := m.Images.GetBestCandidate()

		model.Width = best.Width
		model.Height = best.Height
		model.Url = best.Url
		model.Type = instagram.MediaTypeImage
	// video
	case MediaTypeVideo:
		best := m.Videos.GetBestCandidate()

		model.Width = best.Width
		model.Height = best.Height
		model.Url = best.Url
		model.Type = instagram.MediaTypeVideo
	// carousel
	// case 8:
	default:
		model.Type = instagram.MediaTypeUndefined
	}

	return model, nil
}

// GetBest returns best quality image or video.
// Arguments can be []Video or []Candidate
// func MediaGetBest(obj interface{}) string {
// 	type bestMedia struct {
// 		w, h int
// 		url  string
// 	}
//
// 	m := bestMedia{}
//
// 	switch t := obj.(type) {
// 	// getting best video
// 	case []Video:
// 		for _, video := range t {
// 			if m.w < video.Width && video.Height > m.h && video.URL != "" {
// 				m.w = video.Width
// 				m.h = video.Height
// 				m.url = video.URL
// 			}
// 		}
// 		// getting best image
// 	case []Candidate:
// 		for _, image := range t {
// 			if m.w < image.Width && image.Height > m.h && image.URL != "" {
// 				m.w = image.Width
// 				m.h = image.Height
// 				m.url = image.URL
// 			}
// 		}
// 	}
// 	return m.url
// }
