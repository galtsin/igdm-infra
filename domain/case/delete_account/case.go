package delete_account

import (
	"fmt"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/case/add_activity_log"
	"channels-instagram-dm/domain/model"
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
	runtimeContext.Logger().Info("[delete_account] Case run", nil)

	err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[delete_account] Case err [%s]", err), nil)
		return err
	}

	return nil
}

func run(runtimeContext domain.RuntimeContext, req Request) error {
	if err := validate(req); err != nil {
		return domain.NewErrorInvalidArgument(err.Error())
	}

	accountRepository := runtimeContext.Repository().AccountRepository()
	credentialRepository := runtimeContext.Repository().CredentialsRepository()

	account, err := accountRepository.WhereExternalID(req.ExternalID)
	if err != nil {
		return err
	}

	if account.State != model.AccountStateSuspend {
		return fmt.Errorf("Account is not suspended")
	}

	credentials, err := credentialRepository.WhereExternalID(req.ExternalID)
	if err != nil {
		return err
	}

	account.State = model.AccountStateSuspend
	account.StateReason = model.AccountStateReasonMarkedAsDeleted

	if _, err = accountRepository.Store(account); err != nil {
		return err
	}

	if err := credentialRepository.Delete(credentials.ID); err != nil {
		return err
	}

	_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
		AccountID: account.ID,
		Log:       "Account was deleted",
	})

	runtimeContext.EventBus().PublishAccountDeleted(domain.EventAccountDeleted{Account: account})

	return nil
}
