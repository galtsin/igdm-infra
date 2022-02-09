package model

import (
	"fmt"

	"channels-instagram-dm/domain/model/channels"
	"channels-instagram-dm/domain/model/instagram"
)

func GetChannelsMessagePayload(m channels.Message) interface{} {
	switch m.Type {
	case channels.MessageTypeText:
		return MessageText{Text: m.Text}
	default:
		return MessageUndefined{
			Text: fmt.Sprintf("unsupported message type %s", m.Type),
		}
	}
}

func GetInstagramMessagePayload(i instagram.ThreadItem) interface{} {
	switch i.Type {
	case instagram.MessageTypeText:
		return MessageText{
			Text: string(i.Text),
		}
	case instagram.MessageTypeLike:
		return MessageLike{
			Like: string(i.Text),
		}
	case instagram.MessageTypeLink:
		return MessageLink{
			Url:             i.Link.Url,
			Title:           i.Link.Title,
			Summary:         i.Link.Summary,
			ImagePreviewUrl: i.Link.ImagePreviewUrl,
		}
	case instagram.MessageTypeActionLog:
		return MessageActionLog{
			Text: string(i.Text),
		}
	case instagram.MessageTypeMedia:
		switch i.Media.Type {
		case instagram.MediaTypeImage:
			return MessageMediaImage{
				MessageMedia: MessageMedia{
					ID:  i.Media.ID,
					Url: i.Media.Url,
				},
				Width:  i.Media.Width,
				Height: i.Media.Height,
			}
		case instagram.MediaTypeVideo:
			return MessageMediaVideo{
				MessageMedia: MessageMedia{
					ID:  i.Media.ID,
					Url: i.Media.Url,
				},
				Width:  i.Media.Width,
				Height: i.Media.Height,
			}
		case instagram.MediaTypeVisualImage:
			return MessageMediaVisualImage{
				MessageMedia: MessageMedia{
					ID:  i.Media.ID,
					Url: i.Media.Url,
				},
				Width:  i.Media.Width,
				Height: i.Media.Height,
			}
		case instagram.MediaTypeVisualVideo:
			return MessageMediaVisualVideo{
				MessageMedia: MessageMedia{
					ID:  i.Media.ID,
					Url: i.Media.Url,
				},
				Width:  i.Media.Width,
				Height: i.Media.Height,
			}
		case instagram.MediaTypeAnimated:
			return MessageMediaAnimated{
				MessageMedia: MessageMedia{
					ID:  i.Media.ID,
					Url: i.Media.Url,
				},
				Width:  i.Media.Width,
				Height: i.Media.Height,
			}
		case instagram.MediaTypeVoice:
			return MessageMediaVoice{
				MessageMedia: MessageMedia{
					ID:  i.Media.ID,
					Url: i.Media.Url,
				},
			}
		default:
			return MessageUndefined{
				Text: fmt.Sprintf("unsupported message media type "),
			}
		}
	default:
		return MessageUndefined{
			Text: fmt.Sprintf("unsupported message type "),
		}
	}
}
