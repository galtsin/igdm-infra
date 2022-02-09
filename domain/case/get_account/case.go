package get_account

import (
	"fmt"

	"channels-instagram-dm/domain"
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
	runtimeContext.Logger().Info("[get_account] Case run", nil)

	resp, err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[get_account] Case err [%s]", err), nil)
		return resp, err
	}

	return resp, nil
}

func run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	resp := Response{}

	if err := validate(req); err != nil {
		return resp, domain.NewErrorInvalidArgument(err.Error())
	}

	accountRepository := runtimeContext.Repository().AccountRepository()

	account, err := accountRepository.WhereExternalID(req.ExternalID)
	if err != nil {
		return resp, err
	}

	resp.Account = account

	return resp, nil

}
