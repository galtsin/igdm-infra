package jsonapi

import (
	"encoding/json"

	"channels-instagram-dm/domain"
)

type ErrorPresenter interface {
	Marshal(err domain.BaseError) []byte
}

type errorPresenter struct{}

type Error struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

func NewErrorPresenter() ErrorPresenter {
	return &errorPresenter{}
}

func (p *errorPresenter) Marshal(errBase domain.BaseError) []byte {
	result := struct {
		Error Error `json:"error"`
	}{
		Error: Error{
			Code:  errBase.Code(),
			Title: errBase.Error(),
		},
	}

	bs, err := json.Marshal(result)
	if err != nil {
		return []byte(`{"error":{"title":"Unable to format error"}}`)
	}

	return bs
}
