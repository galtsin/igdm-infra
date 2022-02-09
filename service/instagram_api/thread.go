package instagram_api

import (
	"context"
	"fmt"
	"time"

	"channels-instagram-dm/domain/model/instagram"
)

func (s *service) DirectThread(threadID, cursor string) (resp instagram.ThreadWithItems, err error) {
	request := NewRequest("direct@thread")
	request.Params = struct {
		ThreadID string `json:"thread_id"`
		Cursor   string `json:"cursor,omitempty"`
	}{
		ThreadID: threadID,
		Cursor:   cursor,
	}

	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Minute)
	defer cancel()

	ch, err := s.send(ctx, request.ID, request, new(Thread))
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

		if val, ok := response.Result.(*Thread); ok {
			return val.toModel()
		}

		return resp, fmt.Errorf("")
	}
}
