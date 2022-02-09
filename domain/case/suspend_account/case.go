package suspend_account

import (
	"fmt"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/case/add_activity_log"
	"channels-instagram-dm/domain/model"
)

type Request struct {
	ExternalID string
	StopReason string
}

func validate(req Request) error {
	if req.ExternalID == "" {
		return fmt.Errorf("ExternalID should not be empty")
	}

	if req.StopReason == "" {
		return fmt.Errorf("StopReason should not be empty")
	}

	return nil
}

func Run(runtimeContext domain.RuntimeContext, req Request) error {
	runtimeContext.Logger().Info("[suspend_account] Case run", nil)

	err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[suspend_account] Case err [%s]", err), nil)
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

	// Дальнейшее изменение suspend недопустимо
	if account.State == model.AccountStateSuspend {
		runtimeContext.Logger().Info(fmt.Sprintf("[suspend_account] Account [%s] is suspended", req.ExternalID), nil)
		return nil
	}

	account.State = model.AccountStateSuspend
	account.StateReason = req.StopReason

	if _, err = accountRepository.Store(account); err != nil {
		return err
	}

	// Activity log
	log := ""

	switch account.StateReason {
	case model.AccountStateReasonServiceStopped:
		log = "stopped by service"
	case model.AccountStateReasonNoLoggedIn:
		log = "no logged in"
	case model.AccountStateReasonPermanentError:
		log = "permanent error"
	case model.AccountStateReasonChallenge:
		log = "challenge required"
	case model.AccountStateReasonMarkedAsDeleted:
		log = "marked as deleted"
	default:
		log = account.StateReason
	}

	_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
		AccountID: account.ID,
		Log:       "Account was suspended due " + log,
	})

	runtimeContext.EventBus().PublishAccountSuspended(domain.EventAccountSuspended{
		Account: account,
	})

	return nil

}
