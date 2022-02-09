package get_all_accounts

import (
	"fmt"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
)

type Request struct {
}

type Response struct {
	Accounts []model.Account
}

func Run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	runtimeContext.Logger().Info("[get_all_accounts] Case run", nil)

	resp, err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[get_all_accounts] Case err [%s]", err), nil)
		return resp, err
	}

	return resp, nil
}

func run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	resp := Response{}
	accountRepository := runtimeContext.Repository().AccountRepository()

	accounts, err := accountRepository.All()
	if err != nil {
		return resp, err
	}

	resp.Accounts = accounts

	return resp, nil

}
