package add_account

import (
	"fmt"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/case/add_activity_log"
	"channels-instagram-dm/domain/model"
)

type Request struct {
	ExternalID string
}

type Response struct {
	Account model.Account
}

func validate(req Request) error {
	if req.ExternalID == "" {
		return fmt.Errorf("ExternalID should not be empty")
	}

	return nil
}

func Run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	runtimeContext.Logger().Info("[add_account] Case run", nil)

	resp, err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[add_account] Case err [%s]", err), nil)
		return resp, err
	}

	return resp, nil
}

func run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	resp := Response{}

	if err := validate(req); err != nil {
		return resp, domain.NewErrorInvalidArgument(err.Error())
	}

	credentials, err := runtimeContext.Repository().CredentialsRepository().WhereExternalID(req.ExternalID)
	if err != nil {
		return resp, err
	}

	account := model.NewAccount(req.ExternalID, credentials.Username)
	account.State = model.AccountStateActive

	account, err = runtimeContext.Repository().AccountRepository().Store(account)
	if err != nil {
		return resp, err
	}

	_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
		AccountID: account.ID,
		Log:       "Account was created",
	})

	resp.Account = account

	runtimeContext.EventBus().PublishAccountCreated(domain.EventAccountCreated{Account: account})

	return resp, nil
}
