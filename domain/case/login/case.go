package login

import (
	"errors"
	"fmt"
	"net/url"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/case/add_activity_log"
	"channels-instagram-dm/domain/model"
	"channels-instagram-dm/domain/model/instagram"
)

type Request struct {
	// Автоматический релогин - инструкция, чтобы использовать текущие учетные данные без изменений
	AutoLogin bool
	Login     instagram.Login
}

type Response struct {
	Login instagram.Login
}

func validate(req Request) error {
	if req.Login.ExternalID == "" {
		return fmt.Errorf("ExternalID should not be empty")
	}

	if req.Login.Required.Case == instagram.RequiredStep2F {
		if req.Login.Required.Options.Identifier == "" {
			return fmt.Errorf("2f identifier should not be empty")
		}

		if req.Login.Required.Options.Method == "" {
			return fmt.Errorf("2f method should not be empty")
		}

		if req.Login.Required.Options.Code == "" {
			return fmt.Errorf("2f code should not be empty")
		}
	}

	if req.Login.Required.Case == instagram.RequiredStepChallenge {
		if req.Login.Required.Options.Step == "" {
			return fmt.Errorf("Challenge step should not be empty")
		}

		if req.Login.Required.Options.CheckpointUrl == "" {
			return fmt.Errorf("Challenge checkpoint url  should not be empty")
		}

		if req.Login.Required.Options.Code == "" {
			return fmt.Errorf("Challenge code should not be empty")
		}
	}

	return nil
}

func Run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	runtimeContext.Logger().Info("[login] Case run", nil)

	resp, err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[login] Case err [%s]", err), nil)
		return resp, err
	}

	return resp, nil
}

func run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	resp := Response{
		Login: req.Login,
	}

	if err := validate(req); err != nil {
		return resp, domain.NewErrorInvalidArgument(err.Error())
	}

	credentialsRepository := runtimeContext.Repository().CredentialsRepository()

	credentials, err := credentialsRepository.WhereExternalID(req.Login.ExternalID)

	if err != nil {
		if !errors.Is(err, domain.ErrorNotFound) {
			return resp, err
		}

		// Нет записи о credentials
		// Некоторые поля обязательны для заполнения
		if req.AutoLogin {
			return resp, domain.NewErrorInvalidArgument("Credentials does not exist for autologin")
		}

		if req.Login.Credentials.Username == "" {
			return resp, domain.NewErrorInvalidArgument("Username should not be empty")
		}

		if req.Login.Credentials.Password == "" {
			return resp, domain.NewErrorInvalidArgument("Password should not be empty")
		}

		if req.Login.Credentials.Proxy != "" {
			if _, err := url.ParseRequestURI(req.Login.Credentials.Proxy); err != nil {
				return resp, domain.NewErrorInvalidArgument(fmt.Sprintf("Proxy is invalid. %s", err))
			}
		}

		credentials = model.Credentials{
			ExternalID: req.Login.ExternalID,
			Username:   req.Login.Credentials.Username,
			Password:   req.Login.Credentials.Password,
			Proxy:      req.Login.Credentials.Proxy,
		}

	}

	if !req.AutoLogin && credentials.ID != "" {
		// Пароль может быть пустым, в таком случае его не изменяем
		if req.Login.Credentials.Password != "" {
			credentials.Password = req.Login.Credentials.Password
		}

		// Прокси может быть пустым и это валидный случай
		if req.Login.Credentials.Proxy != "" {
			if _, err := url.ParseRequestURI(req.Login.Credentials.Proxy); err != nil {
				return resp, domain.NewErrorInvalidArgument(fmt.Sprintf("Proxy is invalid. %s", err))
			}
		}

		credentials.Proxy = req.Login.Credentials.Proxy
	}

	login := instagram.Login{
		ExternalID: credentials.ExternalID,
		Credentials: instagram.Credentials{
			Username: credentials.Username,
			Password: credentials.Password,
			Proxy:    credentials.Proxy,
		},
		Required: req.Login.Required,
	}

	api, err := runtimeContext.Service().InstagramAPI(login.Credentials.Username)
	if err != nil {
		return resp, err
	}

	switch req.Login.Required.Case {
	case instagram.RequiredStep2F:
		if err := api.Login2F(login.Credentials, login.Required); err != nil {
			return resp, err
		}

		resp.Login.Required = instagram.Required{
			Case: instagram.RequiredStepNone,
		}
	case instagram.RequiredStepChallenge:
		if err := api.Challenge(login.Required); err != nil {
			return resp, err
		}

		resp.Login.Required = instagram.Required{
			Case: instagram.RequiredStepNone,
		}
	default:
		loginRequired, err := api.Login(login.Credentials)
		if err != nil {
			return resp, err
		}

		resp.Login.Required = loginRequired
	}

	_, err = credentialsRepository.Store(credentials)
	if err != nil {
		return resp, err
	}

	// Activity log
	account, err := runtimeContext.Repository().AccountRepository().WhereExternalID(credentials.ExternalID)

	if err == nil && account.ID != "" {
		switch resp.Login.Required.Case {
		case instagram.RequiredStepNone:
			if req.Login.Credentials.Password != "" {
				_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
					AccountID: account.ID,
					Log:       "Credentials was changed",
				})
			}

			if req.AutoLogin {
				// Возможно такое не стоит сохранять в активности
				_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
					AccountID: account.ID,
					Log:       "Auto login success",
				})
			} else {
				_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
					AccountID: account.ID,
					Log:       "Login success",
				})
			}

		case instagram.RequiredStep2F:
			_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
				AccountID: account.ID,
				Log:       "Login require two factor",
			})

		case instagram.RequiredStepChallenge:
			_, _ = add_activity_log.Run(runtimeContext, add_activity_log.Request{
				AccountID: account.ID,
				Log:       "Login require challenge",
			})
		}
	}

	return resp, nil
}
