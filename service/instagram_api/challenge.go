package instagram_api

import (
	"context"
	"fmt"
	"time"

	"channels-instagram-dm/domain/model/instagram"
)

type ChallengeParams struct {
	Step          string `json:"step"`
	CheckpointUrl string `json:"checkpoint_url"`
	Code          string `json:"code"`
}

func (s *service) Challenge(required instagram.Required) error {
	if required.Case != instagram.RequiredStepChallenge {
		return fmt.Errorf("Expecting challenge step. Occurred is %s", string(required.Case))
	}

	request := NewRequest("auth@challenge")
	request.Params = ChallengeParams{
		Step:          required.Options.Step,
		CheckpointUrl: required.Options.CheckpointUrl,
		Code:          required.Options.Code,
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
			return newErrorChallengeFailed(response.Error)
		}

		return nil
	}
}
