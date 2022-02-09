package sync

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/case/login"
	"channels-instagram-dm/domain/case/suspend_account"
	"channels-instagram-dm/domain/model"
	"channels-instagram-dm/domain/model/instagram"
	sync_instagram "channels-instagram-dm/sync/instagram"
	sync_undelivered_message "channels-instagram-dm/sync/undelivered_message"
	utility_clean_activity_log "channels-instagram-dm/sync/utility"
)

// Внимание: Остановка происходит через прерывание контекста
func startAccount(runtimeContext domain.RuntimeContext, account model.Account) chan struct{} {
	done := make(chan struct{}, 1)

	go func() {
		defer func() {
			done <- struct{}{}
		}()

		if err := tryLogin(runtimeContext, account); err != nil {
			runtimeContext.Logger().Error(fmt.Sprintf("Failed to IG Login. %s", err), nil)

			runtimeContext.EventBus().PublishSuspendAccount(domain.EventSuspendAccount{
				Reason:  model.AccountStateReasonNoLoggedIn,
				Account: account,
			})

			return
		}

		// Канал транспорта сообщений Instagram-Channels
		chTransfer := make(chan model.MessagesBatch, 0)

		defer func() {
			close(chTransfer)
		}()

		// Синхронизация горутин
		wg := &sync.WaitGroup{}

		defer func() {
			wg.Wait()
		}()

		var launch error

		defer func() {
			if launch != nil {
				runtimeContext.EventBus().PublishSuspendAccount(domain.EventSuspendAccount{
					Reason:  model.AccountStateReasonPermanentError,
					Account: account,
				})

				runtimeContext.Logger().Error(fmt.Sprintf("Failed to launch account. %s", launch), nil)
			}
		}()

		producer := runtimeContext.MQ().Producer()
		consumer := runtimeContext.MQ().Consumer()

		if producer == nil {
			launch = errors.New("Producer is nil")
			return
		}

		if consumer == nil {
			launch = errors.New("Consumer is nil")
			return
		}

		defer func() {
			if err := consumer.Close(); err != nil {
				runtimeContext.Logger().Error(fmt.Sprintf("Failed to close consumer. %s", err), nil)
			}
		}()

		if err := sync_instagram.Listen(runtimeContext.WithLogger(runtimeContext.Logger().Copy("INSTAGRAM")), wg, account); err != nil {
			launch = fmt.Errorf("Unable to start sync instagram. %s ", err)
			return
		}

		if err := sync_undelivered_message.Listen(runtimeContext.WithLogger(runtimeContext.Logger().Copy("UNDELIVERED")), wg, chTransfer, account); err != nil {
			launch = fmt.Errorf("Unable to start sync undelivered messages. %s ", err)
			return
		}

		if err := utility_clean_activity_log.Listen(runtimeContext.WithLogger(runtimeContext.Logger().Copy("CLEAN")), wg, account); err != nil {
			launch = fmt.Errorf("Unable to start utility clean_activity_log. %s ", err)
			return
		}

		if launch != nil {
			return
		}

		chSubscribeOnLoginAccount := runtimeContext.EventBus().SubscribeOnLoginAccount(func(event interface{}) bool {
			if event, ok := event.(domain.EventLoginAccount); ok {
				return event.Account.ID == account.ID
			}

			return false
		})

		// Пытаемся переавторизоваться каждые X
		tickerDefaultDuration := 8 * time.Hour
		tryLoginAttempts := 0

		ticker := time.NewTicker(tickerDefaultDuration)

		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case <-runtimeContext.Context().Done():
				return
			case <-chSubscribeOnLoginAccount:
				runtimeContext.Logger().Debug("Event SubscribeOnLoginAccount", nil)
				// Лаг на отложенное выполнение, чтобы избежать каскадной обработки событий, схлопывая серию в одно событие
				// Сдвигаем событие на величину временного окна Х
				ticker.Reset(15 * time.Second)
			case <-ticker.C:
				ticker.Reset(tickerDefaultDuration)

				if tryLoginAttempts == 5 {
					runtimeContext.EventBus().PublishSuspendAccount(domain.EventSuspendAccount{
						Reason:  model.AccountStateReasonPermanentError,
						Account: account,
					})

					return
				}

				err := tryLogin(runtimeContext, account)

				if err == nil {
					tryLoginAttempts = 0
					ticker.Reset(tickerDefaultDuration)
					continue
				}

				if errors.Is(err, domain.ErrorNoLoggedIn) || errors.Is(err, domain.ErrorInvalidCredentials) {
					runtimeContext.EventBus().PublishSuspendAccount(domain.EventSuspendAccount{
						Reason:  model.AccountStateReasonNoLoggedIn,
						Account: account,
					})

					return
				}

				tryLoginAttempts++
				ticker.Reset(5 * time.Minute)
			}
		}
	}()

	return done
}

func stopAccount(runtimeContext domain.RuntimeContext, account model.Account, reason string) error {
	return suspend_account.Run(runtimeContext, suspend_account.Request{
		ExternalID: account.ExternalID,
		StopReason: reason,
	})
}

func tryLogin(runtimeContext domain.RuntimeContext, account model.Account) error {
	resp, err := login.Run(runtimeContext, login.Request{
		AutoLogin: true,
		Login: instagram.Login{
			ExternalID: account.ExternalID,
		}})

	if err != nil {
		return err
	}

	if resp.Login.Required.Case != instagram.RequiredStepNone {
		return domain.NewErrorNoLoggedIn(fmt.Sprintf("Required is %s", resp.Login.Required.Case))
	}

	return nil
}
