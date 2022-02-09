package types

import (
	"channels-instagram-dm/domain/model/instagram"
)

type Like string

func (m Like) ToModel() (instagram.Text, error) {
	return instagram.Text(m), nil
}
