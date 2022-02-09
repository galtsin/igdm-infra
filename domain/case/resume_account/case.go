package resume_account

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
	runtimeContext.Logger().Info("[resume_account] Case run", nil)

	err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[resume_account] Case err [%s]", err), nil)
		return err
	}

	return nil
}

func run(runtimeContext domain.RuntimeContext, req Request) error {
	if err := validate(req); err != nil {
		return domain.NewErrorInvalidArgument(err.Error())
	}

	accountRepository := runtimeContext.Repository().AccountRepository()
	account, err := accountRepository.WhereExternalID(req.ExternalID)
	if err != nil {
		return err
	}

	if account.State == model.AccountStateActive {
		runtimeContext.Logger().Info(fmt.Sprintf("[resume_account] Account [%s] is active", req.ExternalID), nil)
		return nil
	}

	if account.StateReason == model.AccountStateReasonMarkedAsDeleted {
		return fmt.Errorf("Account marked as deleted ")
	}

	account.State = model.AccountStateActive
	account.StateReason = ""

	acc, err := accountRepository.Store(account)
	if err != nil {
		return err
	}

	_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
		AccountID: account.ID,
		Log:       "Account was resumed",
	})

	runtimeContext.EventBus().PublishAccountResumed(domain.EventAccountResumed{
		Account: acc,
	})

	return nil
}
