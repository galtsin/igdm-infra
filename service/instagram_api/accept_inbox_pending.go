package instagram_api

import (
	"context"
	"fmt"
	"time"
)

func (s *service) DirectAcceptInboxPending(threadIDs []string) error {
	request := NewRequest("direct@accept_inbox_pending")
	request.Params = struct {
		Threads []string `json:"threads"`
	}{
		Threads: threadIDs,
	}

	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Minute)
	defer cancel()

	ch, err := s.send(ctx, request.ID, request, new([]string))
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("Stopped by timeout %w", ctx.Err())
	case response, ok := <-ch:
		if !ok {
			return fmt.Errorf("Channel is closed")
		}

		if response.Error != "" {
			return newError(response.Error)
		}

		if _, ok := response.Result.(*[]string); ok {
			return nil
		}

		return fmt.Errorf("Unhandled error")
	}
}
