package utility

import (
	"fmt"
	"sync"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/case/clean_activity_log"
	"channels-instagram-dm/domain/model"
)

func Listen(runtimeContext domain.RuntimeContext, wg *sync.WaitGroup, account model.Account) error {
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()

		ticker := time.NewTicker(12 * time.Hour)

		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case <-runtimeContext.Context().Done():
				runtimeContext.Logger().Debug("Context was closed", nil)
				return
			case <-ticker.C:
			}

			err := clean_activity_log.Run(runtimeContext, clean_activity_log.Request{
				AccountID: account.ID,
			})

			if err != nil {
				runtimeContext.Logger().Error(fmt.Sprintf("%s", err), nil)
			}
		}
	}()

	return nil
}
