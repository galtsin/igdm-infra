package instagram_api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model/instagram"
)

func (s *service) Login(credentials instagram.Credentials) (instagram.Required, error) {
	request := NewRequest("auth@login")
	request.Params = struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Proxy    string `json:"proxy"`
	}{
		Username: credentials.Username,
		Password: credentials.Password,
		Proxy:    credentials.Proxy,
	}

	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Minute)
	defer cancel()

	required := instagram.Required{}

	ch, err := s.send(ctx, request.ID, request, new(LoginRequired))
	if err != nil {
		return required, err
	}

	select {
	case <-ctx.Done():
		return required, fmt.Errorf("Stopped by timeout %w", ctx.Err())
	case response, ok := <-ch:
		if !ok {
			return required, fmt.Errorf("Channel was closed")
		}

		if response.Error != "" {
			if errors.Is(newError(response.Error), domain.ErrorInvalidCredentials) {
				return required, newError(response.Error)
			}

			return required, newErrorLoginFailed(response.Error)
		}

		if val, ok := response.Result.(*LoginRequired); ok {
			return val.toModel(), nil
		}

		return required, fmt.Errorf("Unhandled error")
	}
}
