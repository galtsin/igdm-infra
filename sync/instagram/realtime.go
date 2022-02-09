package instagram

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
	"channels-instagram-dm/domain/model/instagram"
)

func listenRealtime(runtimeContext domain.RuntimeContext, wg *sync.WaitGroup, chUpdates chan instagram.RealtimeUpdate, account model.Account) error {
	wg.Add(1)

	// Обрабатываем разрывы соединения
	go func() {
		defer func() {
			// close(chUpdates)
			wg.Done()
		}()

		ticker := time.NewTicker(3 * time.Minute)

		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case <-runtimeContext.Context().Done():
				return
			case <-ticker.C:
			}

			// Переподключаемся к слотам, если были прерывания с библиотекой
			// Пытаемся перелогиниться перед обновление
			instagramAPI, err := runtimeContext.Service().InstagramAPI(account.Username)
			if err != nil {
				runtimeContext.Logger().Error(fmt.Sprintf("Failed get Instagram. %s", err), nil)
				continue
			}

			if instagramAPI == nil {
				continue
			}

			if instagramAPI.IsClosed() {
				continue
			}

			closed, err := instagramAPI.ListenThreadUpdates(chUpdates)
			if err != nil {
				runtimeContext.Logger().Error(fmt.Sprintf("Failed to listen thread updates. %s", err), nil)
				continue
			}

			if closed == nil {
				continue
			}

			runtimeContext.Logger().Debug("Listening", nil)

			select {
			case <-runtimeContext.Context().Done():
				runtimeContext.Logger().Debug("Context was closed", nil)
				return

			case err := <-closed:
				if err == nil {
					runtimeContext.Logger().Debug("Channel was closed.", nil)
					continue
				}

				if errors.Is(err, domain.ErrorNoLoggedIn) {
					runtimeContext.EventBus().PublishLoginAccount(domain.EventLoginAccount{
						Account: account,
					})
				}

				runtimeContext.Logger().Error(fmt.Sprintf("Channel was closed. %s", err), nil)
				continue
			}
		}
	}()

	return nil
}

func handleRealtime(runtimeContext domain.RuntimeContext, account model.Account, realtimeUpdate instagram.RealtimeUpdate) error {
	return nil
}
