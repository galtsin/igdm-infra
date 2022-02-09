package instagram_api

import (
	"context"
	"fmt"
	"time"
)

func (s *service) Logout() error {
	request := NewRequest("auth@logout")

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
			return newError(response.Error)
		}

		return nil
	}
}
