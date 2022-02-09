package logout

import (
	"errors"
	"fmt"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/case/add_activity_log"
)

type Request struct {
	ExternalID string
}

func validate(req Request) error {
	if req.ExternalID == "" {
		return fmt.Errorf("ExternalID should not be empty")
	}

	return nil
}

func Run(runtimeContext domain.RuntimeContext, req Request) error {
	runtimeContext.Logger().Info("[logout] Case run", nil)

	err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[logout] Case err [%s]", err), nil)
		return err
	}

	return nil
}

func run(runtimeContext domain.RuntimeContext, req Request) error {
	if err := validate(req); err != nil {
		return err
	}

	credentials, err := runtimeContext.Repository().CredentialsRepository().WhereExternalID(req.ExternalID)
	if err != nil {
		return err
	}

	api, err := runtimeContext.Service().InstagramAPI(credentials.Username)
	if err != nil {
		return err
	}

	if err := api.Logout(); err != nil {
		if errors.Is(err, domain.ErrorNoLoggedIn) {
			return nil
		}

		return err
	}

	// Аккаунт может и не существовать
	account, err := runtimeContext.Repository().AccountRepository().WhereExternalID(req.ExternalID)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			return nil
		}

		return err
	}

	_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
		AccountID: account.ID,
		Log:       "Logout completed",
	})

	runtimeContext.EventBus().PublishAccountLogout(domain.EventAccountLogout{Account: account})

	return nil
}
