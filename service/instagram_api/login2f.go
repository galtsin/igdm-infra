package instagram_api

import (
	"context"
	"fmt"
	"time"

	"channels-instagram-dm/domain/model/instagram"
)

type Login2FParams struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Identifier string `json:"identifier"`
	Method     string `json:"method"`
	Code       string `json:"code"`
}

func (s *service) Login2F(credentials instagram.Credentials, required instagram.Required) error {
	if required.Case != instagram.RequiredStep2F {
		return fmt.Errorf("Expecting 2f step. Occurred is %s", string(required.Case))
	}

	request := NewRequest("auth@login2f")
	request.Params = Login2FParams{
		Username:   credentials.Username,
		Password:   credentials.Password,
		Identifier: required.Options.Identifier,
		Method:     required.Options.Method,
		Code:       required.Options.Code,
	}

	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Minute)
	defer cancel()

	ch, err := s.send(ctx, request.ID, request, nil)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("Stopped by timeout %w", ctx.Err())
	case response, ok := <-ch:
		if !ok {
			return fmt.Errorf("Channel was closed")
		}

		if response.Error != "" {
			return newErrorLoginFailed(response.Error)
		}

		return nil
	}
}
