package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/case/add_account"
	"channels-instagram-dm/domain/case/delete_account"
	"channels-instagram-dm/domain/case/get_account"
	"channels-instagram-dm/domain/case/get_activity_log"
	"channels-instagram-dm/domain/case/get_all_accounts"
	"channels-instagram-dm/domain/case/login"
	"channels-instagram-dm/domain/case/logout"
	"channels-instagram-dm/domain/case/resume_account"
	"channels-instagram-dm/domain/case/suspend_account"
	"channels-instagram-dm/presenter/jsonapi"

	"github.com/gorilla/mux"
)

func HealthCheck(ctx domain.RuntimeContext, req *http.Request) ([]byte, error) {
	return []byte(`{"status":true, "version":1}`), nil
}

func Slots(ctx domain.RuntimeContext, req *http.Request) ([]byte, error) {
	presenter := jsonapi.NewSlotPresenter()
	result, err := presenter.MarshalList(ctx.Service().Slots())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func RefreshSlots(ctx domain.RuntimeContext, req *http.Request) ([]byte, error) {
	err := ctx.Service().RefreshSlots()
	if err != nil {
		return nil, err
	}

	presenter := jsonapi.NewSlotPresenter()
	result, err := presenter.MarshalList(ctx.Service().Slots())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetAllAccounts(runtimeContext domain.RuntimeContext, req *http.Request) ([]byte, error) {
	resp, err := get_all_accounts.Run(runtimeContext, get_all_accounts.Request{})
	if err != nil {
		return nil, err
	}

	presenter := jsonapi.NewAccountPresenter()
	result, err := presenter.MarshalList(resp.Accounts)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetAccount(runtimeContext domain.RuntimeContext, req *http.Request) ([]byte, error) {
	vars := mux.Vars(req)

	resp, err := get_account.Run(runtimeContext, get_account.Request{
		ExternalID: vars["external_id"],
	})
	if err != nil {
		return nil, err
	}

	presenter := jsonapi.NewAccountPresenter()
	result, err := presenter.Marshal(resp.Account)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func AddAccount(runtimeContext domain.RuntimeContext, req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	account, err := jsonapi.NewAccountPresenter().Unmarshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := add_account.Run(runtimeContext, add_account.Request{
		ExternalID: account.ExternalID,
	})
	if err != nil {
		return nil, err
	}

	return jsonapi.NewAccountPresenter().
		Marshal(resp.Account)
}

func Login(runtimeContext domain.RuntimeContext, req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	presenter := jsonapi.NewLoginPresenter()
	lg, err := presenter.Unmarshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := login.Run(runtimeContext, login.Request{
		Login: lg,
	})
	if err != nil {
		return nil, err
	}

	if resp.Login.Required.Case != "" {
		return presenter.MarshalRequired(resp.Login)
	}

	return presenter.Marshal(resp.Login)
}

func Logout(runtimeContext domain.RuntimeContext, req *http.Request) ([]byte, error) {
	vars := mux.Vars(req)

	err := logout.Run(runtimeContext, logout.Request{
		ExternalID: vars["external_id"],
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func ResumeAccount(runtimeContext domain.RuntimeContext, req *http.Request) ([]byte, error) {
	vars := mux.Vars(req)

	err := resume_account.Run(runtimeContext, resume_account.Request{
		ExternalID: vars["external_id"],
	},
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func SuspendAccount(runtimeContext domain.RuntimeContext, req *http.Request) ([]byte, error) {
	vars := mux.Vars(req)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	data := struct {
		Reason string `json:"reason"`
	}{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	err = suspend_account.Run(runtimeContext, suspend_account.Request{
		ExternalID: vars["external_id"],
		StopReason: data.Reason,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func DeleteAccount(runtimeContext domain.RuntimeContext, req *http.Request) ([]byte, error) {
	vars := mux.Vars(req)

	err := delete_account.Run(runtimeContext, delete_account.Request{
		ExternalID: vars["external_id"],
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetActivityLog(runtimeContext domain.RuntimeContext, req *http.Request) ([]byte, error) {
	vars := mux.Vars(req)

	resp, err := get_activity_log.Run(runtimeContext, get_activity_log.Request{
		ExternalID: vars["external_id"],
	})
	if err != nil {
		return nil, err
	}

	presenter := jsonapi.NewActivityLogPresenter()
	result, err := presenter.MarshalList(resp.ActivityLogList)
	if err != nil {
		return nil, err
	}

	return result, nil
}
