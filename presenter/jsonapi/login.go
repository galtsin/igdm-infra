package jsonapi

import (
	"encoding/json"

	"channels-instagram-dm/domain/model/instagram"
)

type LoginPresenter interface {
	Marshal(instagram.Login) ([]byte, error)
	MarshalCredentials(instagram.Login) ([]byte, error)
	MarshalRequired(instagram.Login) ([]byte, error)
	Unmarshal([]byte) (instagram.Login, error)
}

type loginPresenter struct{}

type Login struct {
	Type
	Attributes LoginAttributes `json:"attributes"`
}

type LoginCredentials struct {
	Type
	Attributes Credentials `json:"attributes"`
}

type LoginRequired struct {
	Type
	Attributes Required `json:"attributes"`
}

type LoginAttributes struct {
	Credentials Credentials `json:"credentials"`
	Required    Required    `json:"required"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Proxy    string `json:"proxy"`
}

type Required struct {
	Case    string          `json:"case"`
	Options RequiredOptions `json:"options"`
}

type RequiredOptions struct {
	Identifier    string `json:"identifier"`
	Step          string `json:"step"`
	CheckpointUrl string `json:"checkpoint_url"`
	Method        string `json:"method"`
	Code          string `json:"code"`
}

func NewLoginPresenter() LoginPresenter {
	return &loginPresenter{}
}

func (p *loginPresenter) Unmarshal(data []byte) (login instagram.Login, err error) {
	result := struct {
		Data Login `json:"data"`
	}{}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	return result.Data.toModel(), nil
}

func (p *loginPresenter) Marshal(login instagram.Login) ([]byte, error) {
	l := Login{}
	l.fromModel(login)

	result := struct {
		Data Login `json:"data"`
	}{
		Data: l,
	}

	return json.Marshal(result)
}

func (p *loginPresenter) MarshalCredentials(login instagram.Login) ([]byte, error) {
	l := Login{}
	l.fromModel(login)

	result := struct {
		Data LoginCredentials `json:"data"`
	}{
		Data: LoginCredentials{
			Type:       l.Type,
			Attributes: l.Attributes.Credentials,
		},
	}

	return json.Marshal(result)
}

func (p *loginPresenter) MarshalRequired(login instagram.Login) ([]byte, error) {
	l := Login{}
	l.fromModel(login)

	result := struct {
		Data LoginRequired `json:"data"`
	}{
		Data: LoginRequired{
			Type:       l.Type,
			Attributes: l.Attributes.Required,
		},
	}

	return json.Marshal(result)
}

func (l Login) toModel() instagram.Login {
	login := instagram.Login{
		ExternalID: l.Type.ID,
		Credentials: instagram.Credentials{
			Username: l.Attributes.Credentials.Username,
			Password: l.Attributes.Credentials.Password,
			Proxy:    l.Attributes.Credentials.Proxy,
		},
		Required: instagram.Required{},
	}

	switch l.Attributes.Required.Case {
	case string(instagram.RequiredStep2F):
		login.Required.Case = instagram.RequiredStep2F
		login.Required.Options = instagram.RequiredOptions{
			Identifier: l.Attributes.Required.Options.Identifier,
			Method:     l.Attributes.Required.Options.Method,
			Code:       l.Attributes.Required.Options.Code,
		}
	case string(instagram.RequiredStepChallenge):
		login.Required.Case = instagram.RequiredStepChallenge
		login.Required.Options = instagram.RequiredOptions{
			Step:          l.Attributes.Required.Options.Step,
			CheckpointUrl: l.Attributes.Required.Options.CheckpointUrl,
			Code:          l.Attributes.Required.Options.Code,
		}
	default:
		login.Required.Case = instagram.RequiredStepNone
	}

	return login
}

func (l *Login) fromModel(login instagram.Login) {
	l.Type.ID = login.ExternalID
	l.Type.Type = "login"

	l.Attributes.Credentials = Credentials{
		Username: login.Credentials.Username,
		Password: login.Credentials.Password,
		Proxy:    login.Credentials.Proxy,
	}

	l.Attributes.Required = Required{}
	switch login.Required.Case {
	case instagram.RequiredStep2F:
		l.Attributes.Required.Case = string(instagram.RequiredStep2F)
		l.Attributes.Required.Options = RequiredOptions{
			Identifier: login.Required.Options.Identifier,
			Method:     login.Required.Options.Method,
			Code:       login.Required.Options.Code,
		}
	case instagram.RequiredStepChallenge:
		l.Attributes.Required.Case = string(instagram.RequiredStepChallenge)
		l.Attributes.Required.Options = RequiredOptions{
			Step:          login.Required.Options.Step,
			CheckpointUrl: login.Required.Options.CheckpointUrl,
			Code:          login.Required.Options.Code,
		}
	default:
		l.Attributes.Required.Case = string(instagram.RequiredStepNone)
	}
}
