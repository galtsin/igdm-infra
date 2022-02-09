package instagram

import (
	"errors"
	"fmt"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
)

func handleInbox(runtimeContext domain.RuntimeContext, account model.Account) error {
	if err := syncInbox(runtimeContext, account); err != nil {
		if errors.Is(err, domain.ErrorNoLoggedIn) {
			runtimeContext.EventBus().PublishLoginAccount(domain.EventLoginAccount{
				Account: account,
			})
		}

		return fmt.Errorf("%s", err)
	}

	return nil
}

func syncInbox(runtimeContext domain.RuntimeContext, account model.Account) error {
	return nil
}
