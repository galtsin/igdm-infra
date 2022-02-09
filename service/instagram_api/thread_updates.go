package instagram_api

import (
	"fmt"

	"channels-instagram-dm/domain/model/instagram"
)

func (s *service) ListenThreadUpdates(chRT chan instagram.RealtimeUpdate) (chan error, error) {
	request := NewRequest("realtime@start")
	request.Params = struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	ch, err := s.send(s.ctx, request.ID, request, new(RealtimeUpdate))
	if err != nil {
		return nil, err
	}

	closed := make(chan error, 0)

	go func() {
		var errClosed error

		defer func() {
			closed <- errClosed
			close(closed)
		}()

		defer func() {
			s.logger.Info("ListenThreadUpdates: Stopped", nil)
		}()

		defer func() {
			if r := recover(); r != nil {
				s.logger.Critical(fmt.Sprintf("ListenThreadUpdates: Panic %v", r), nil)
			}
		}()

		for {
			select {
			case <-s.ctx.Done():
				return
			case response, ok := <-ch:
				if !ok {
					s.logger.Info("ListenThreadUpdates: Channel was closed", nil)
					return
				}

				if response.Error != "" {
					s.logger.Error(fmt.Sprintf("ListenThreadUpdates: Response with err %s", response.Error), nil)
					errClosed = newError(response.Error)
					return
				}

				val, ok := response.Result.(*RealtimeUpdate)
				if !ok {
					s.logger.Error("ListenThreadUpdates: Mismatched type ", nil)
					continue
				}

				update, err := val.toModel()
				if err != nil {
					s.logger.Error(fmt.Sprintf("ListenThreadUpdates: Failed map to model with err %s", err), nil)
					continue
				}

				select {
				case <-s.ctx.Done():
					return
				case chRT <- update:

				}
			}
		}
	}()

	return closed, nil
}
