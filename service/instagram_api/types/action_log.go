package types

import (
	"channels-instagram-dm/domain/model/instagram"
)

type ActionLog struct {
	Description   string `json:"description"`
	IsReactionLog bool   `json:"is_reaction_log"`
}

func (m ActionLog) ToModel() (instagram.Text, error) {
	return instagram.Text(m.Description), nil
}
