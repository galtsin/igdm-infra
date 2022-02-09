package instagram_api

import (
	"context"
	"fmt"
	"time"

	"channels-instagram-dm/domain/model"
)

func (s *service) DirectSendText(username string, message model.Message) error {
	request := NewRequest("direct@send_text")

	payload, ok := message.Payload.(model.MessageText)
	if !ok {
		return fmt.Errorf("Mismatch message payload type, want text ")
	}

	request.Params = struct {
		Text     string `json:"text"`
		Username string `json:"username"`
	}{
		Text:     payload.Text,
		Username: username,
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
			return newError(response.Error)
		}

		return nil
	}
}

func (s *service) RealtimeSendText(threadID string, message model.Message) error {
	request := NewRequest("realtime@send_text")

	payload, ok := message.Payload.(model.MessageText)
	if !ok {
		return fmt.Errorf("Mismatch message payload type, want text ")
	}

	request.Params = struct {
		Text     string `json:"text"`
		ThreadID string `json:"thread_id"`
	}{
		Text:     payload.Text,
		ThreadID: threadID,
	}

	ctx, cancel := context.WithTimeout(s.ctx, 60*time.Second)
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
