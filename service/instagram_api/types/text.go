package types

import (
	"channels-instagram-dm/domain/model/instagram"
)

type Text string

func (m Text) ToModel() (instagram.Text, error) {
	return instagram.Text(m), nil
}
