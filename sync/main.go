package sync

import (
	"context"
	"fmt"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
)

type terminator struct {
	account model.Account
	cancel  context.CancelFunc
	done    chan struct{}
}

func Run(runtimeContext domain.RuntimeContext) {
	runtimeContext.Syncer().Add()
	terminateMap := make(map[string]terminator)

	go func() {
		defer runtimeContext.Syncer().Remove()

		chSubscribeOnAccountResumed := runtimeContext.EventBus().SubscribeOnAccountResumed(nil)
		chSubscribeOnAccountCreated := runtimeContext.EventBus().SubscribeOnAccountCreated(nil)
		chSubscribeOnAccountSuspended := runtimeContext.EventBus().SubscribeOnAccountSuspended(nil)
		chSubscribeOnSuspendAccount := runtimeContext.EventBus().SubscribeOnSuspendAccount(nil)
		chSubscribeOnAccountLogout := runtimeContext.EventBus().SubscribeOnAccountLogout(nil)

		for {
			select {
			case <-runtimeContext.Context().Done():
				for _, proc := range terminateMap {
					proc.cancel()
					<-proc.done

					logger := runtimeContext.Logger().Copy(proc.account.ExternalID)
					_ = stopAccount(runtimeContext.WithLogger(logger), proc.account, model.AccountStateReasonServiceStopped)
				}

				return

			case e := <-chSubscribeOnAccountResumed:
				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnAccountResumed: Processing with account [%s]", e.Account.ExternalID), nil)

				if _, ok := terminateMap[e.Account.ID]; ok {
					continue
				}

				ctx, cancel := context.WithCancel(runtimeContext.Context())
				logger := runtimeContext.Logger().Copy(e.Account.ExternalID)

				done := startAccount(runtimeContext.WithContext(ctx).WithLogger(logger), e.Account)

				terminateMap[e.Account.ID] = terminator{
					account: e.Account,
					cancel:  cancel,
					done:    done,
				}

				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnAccountResumed: Done with account [%s]", e.Account.ExternalID), nil)

			case e := <-chSubscribeOnAccountCreated:
				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnAccountCreated: Processing with account [%s]", e.Account.ExternalID), nil)

				ctx, cancel := context.WithCancel(runtimeContext.Context())
				logger := runtimeContext.Logger().Copy(e.Account.ExternalID)

				done := startAccount(runtimeContext.WithContext(ctx).WithLogger(logger), e.Account)

				terminateMap[e.Account.ID] = terminator{
					account: e.Account,
					cancel:  cancel,
					done:    done,
				}

				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnAccountCreated: Done with account [%s]", e.Account.ExternalID), nil)

			case e := <-chSubscribeOnAccountSuspended:
				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnAccountSuspended: Processing with account [%s]", e.Account.ExternalID), nil)

				proc, ok := terminateMap[e.Account.ID]
				if !ok {
					continue
				}

				proc.cancel()
				<-proc.done
				delete(terminateMap, e.Account.ID)

				if api, err := runtimeContext.Service().InstagramAPI(e.Account.Username); err != nil {
					runtimeContext.Logger().Error(fmt.Sprintf("Event SubscribeOnAccountSuspended: Failed get InstagramAPI with account [%s]. %s", e.Account.ExternalID, err), nil)
				} else {
					api.Close()
				}

				if err := runtimeContext.Service().RefreshSlots(); err != nil {
					runtimeContext.Logger().Error(fmt.Sprintf("Event SubscribeOnAccountSuspended: Failed RefreshSlots with account [%s]. %s", e.Account.ExternalID, err), nil)
				}

				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnAccountSuspended: Done with account [%s]]", e.Account.ExternalID), nil)

			case e := <-chSubscribeOnSuspendAccount:
				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnSuspendAccount: Processing with account [%s]", e.Account.ExternalID), nil)

				if proc, ok := terminateMap[e.Account.ID]; ok {
					logger := runtimeContext.Logger().Copy(proc.account.ExternalID)

					if err := stopAccount(runtimeContext.WithLogger(logger), proc.account, e.Reason); err != nil {
						runtimeContext.Logger().Error(fmt.Sprintf("Event SubscribeOnSuspendAccount: Failed with account [%s]. %s", e.Account.ExternalID, err), nil)
					}
				}

				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnSuspendAccount: Done with account [%s]", e.Account.ExternalID), nil)

			case e := <-chSubscribeOnAccountLogout:
				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnAccountLogout: Processing with account [%s]", e.Account.ExternalID), nil)

				_, ok := terminateMap[e.Account.ID]

				// Logout, когда аккаунт активен
				if ok {
					runtimeContext.EventBus().PublishSuspendAccount(domain.EventSuspendAccount{
						Reason:  model.AccountStateReasonNoLoggedIn,
						Account: e.Account,
					})
				}

				// Logout, когда аккаунт не добавлен
				if !ok {
					if api, err := runtimeContext.Service().InstagramAPI(e.Account.Username); err != nil {
						runtimeContext.Logger().Error(fmt.Sprintf("Event SubscribeOnAccountLogout: Failed get InstagramAPI with account [%s]. %s", e.Account.ExternalID, err), nil)
					} else {
						api.Close()
					}

					if err := runtimeContext.Service().RefreshSlots(); err != nil {
						runtimeContext.Logger().Error(fmt.Sprintf("Event SubscribeOnAccountLogout: Failed RefreshSlots with account [%s]. %s", e.Account.ExternalID, err), nil)
					}
				}

				runtimeContext.Logger().Info(fmt.Sprintf("Event SubscribeOnAccountLogout: Done with account [%s]", e.Account.ExternalID), nil)

			}
		}
	}()
}
