package instagram_api

import (
	"context"
	"fmt"
	"time"

	"channels-instagram-dm/domain/model/instagram"
)

func (s *service) DirectInbox(cursor string, limit int) (resp instagram.InboxWithThreads, err error) {
	request := NewRequest("direct@inbox")
	request.Params = struct {
		Cursor string `json:"cursor"`
		Limit  int    `json:"limit"`
	}{
		Cursor: cursor,
		Limit:  limit,
	}

	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Minute)
	defer cancel()

	ch, err := s.send(ctx, request.ID, request, new(InboxResponse))
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

		if val, ok := response.Result.(*InboxResponse); ok {
			return val.toModel()
		}

		return resp, fmt.Errorf("Unhandled error")
	}
}
