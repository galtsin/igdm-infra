package instagram_api

import (
	"context"
	"fmt"
	"time"

	"channels-instagram-dm/domain"
)

func (s *service) Discovery() (d domain.DiscoveryRow, err error) {
	request := NewRequest("system@discovery")

	ctx, cancel := context.WithTimeout(s.ctx, 60*time.Second)
	defer cancel()

	ch, err := s.send(ctx, request.ID, request, new(domain.DiscoveryRow))
	if err != nil {
		return
	}

	select {
	case <-ctx.Done():
		return d, fmt.Errorf("Stopped by timeout %w", ctx.Err())
	case response, ok := <-ch:
		if !ok {
			return d, fmt.Errorf("Channel was closed")
		}

		if response.Error != "" {
			return d, newError(response.Error)
		}

		if val, ok := response.Result.(*domain.DiscoveryRow); ok {
			return *val, nil
		}

		return d, fmt.Errorf("Unhandled error")
	}
}
