package instagram

import (
	"fmt"
	"sync"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
	"channels-instagram-dm/domain/model/instagram"
)

func Listen(runtimeContext domain.RuntimeContext, wg *sync.WaitGroup, account model.Account) error {
	chInboxHasChanges := runtimeContext.EventBus().SubscribeOnInboxHasChanges(func(event interface{}) bool {
		if event, ok := event.(domain.EventInboxHasChanges); ok {
			return event.Account.ID == account.ID
		}

		return false
	})

	chRealtimeUpdates := make(chan instagram.RealtimeUpdate, 0)

	realtimeRuntimeContext := runtimeContext.WithLogger(runtimeContext.Logger().Copy("REALTIME"))
	inboxRuntimeContext := runtimeContext.WithLogger(runtimeContext.Logger().Copy("INBOX"))

	err := listenRealtime(realtimeRuntimeContext, wg, chRealtimeUpdates, account)
	if err != nil {
		return err
	}

	// При первом запуске стараемся собрать как можно быстрее
	inboxScheduler := NewScheduler(runtimeContext.Context(), inboxRuntimeContext.Logger())
	inboxScheduler.setRound(jitter(50 * time.Second)) // Джиттер нужен, когда перезапускаем сервис

	wg.Add(1)

	go func() {
		defer func() {
			close(chInboxHasChanges)
			close(chRealtimeUpdates)
			wg.Done()
		}()

		for {
			select {
			case <-runtimeContext.Context().Done():
				inboxRuntimeContext.Logger().Debug("Context was closed", nil)
				return

			case <-chInboxHasChanges:
				inboxRuntimeContext.Logger().Debug(fmt.Sprintf("Event SubscribeOnInboxHasChanges"), nil)

				inboxScheduler.Reset()

				inboxScheduler.Next()

			case <-inboxScheduler.Ticker().C:
				inboxRuntimeContext.Logger().Debug("Scheduler Round", nil)

				if err := handleInbox(inboxRuntimeContext, account); err != nil {
					inboxRuntimeContext.Logger().Error(fmt.Sprintf("Failed to handle. %s", err), nil)

					inboxScheduler.Fail()
				}

				inboxScheduler.Next()

			case update := <-chRealtimeUpdates:
				if err := handleRealtime(realtimeRuntimeContext, account, update); err != nil {
					realtimeRuntimeContext.Logger().Error(fmt.Sprintf("Failed to handle. %s", err), nil)
				}
			}
		}
	}()

	return nil
}
