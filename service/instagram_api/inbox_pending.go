package instagram_api

import (
	"context"
	"fmt"
	"time"

	"channels-instagram-dm/domain/model/instagram"
)

func (s *service) DirectInboxPending(cursor string) (resp []instagram.ThreadWithItems, err error) {
	request := NewRequest("direct@inbox_pending")
	request.Params = struct {
		Cursor string `json:"cursor"`
	}{
		Cursor: cursor,
	}

	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Minute)
	defer cancel()

	ch, err := s.send(ctx, request.ID, request, new([]Thread))
	if err != nil {
		return resp, err
	}

	select {
	case <-ctx.Done():
		return resp, fmt.Errorf("Stopped by timeout %w", ctx.Err())
	case response, ok := <-ch:
		if !ok {
			return resp, fmt.Errorf("Channel is closed")
		}

		if response.Error != "" {
			return resp, newError(response.Error)
		}

		if val, ok := response.Result.(*[]Thread); ok {
			result := make([]instagram.ThreadWithItems, 0, len(*val))

			for _, thread := range *val {
				model, err := thread.toModel()
				if err != nil {
					return resp, err
				}

				result = append(result, model)
			}

			resp = result
			return resp, nil
		}

		return resp, fmt.Errorf("Unhandled error")
	}
}
