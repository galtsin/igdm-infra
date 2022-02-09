package jsonapi

import (
	"encoding/json"
	"time"

	"channels-instagram-dm/domain/model"
)

type AccountPresenter interface {
	Marshal(model.Account) ([]byte, error)
	MarshalList([]model.Account) ([]byte, error)
	Unmarshal([]byte) (model.Account, error)
}

type accountPresenter struct {
	included []*Resource
}

type Account struct {
	Type
	Attributes AccountAttributes `json:"attributes"`
}

type AccountAttributes struct {
	ExternalID  string             `json:"external_id"`
	Username    string             `json:"username"`
	State       model.AccountState `json:"state"`
	StateReason string             `json:"state_reason"`
	CreatedAt   time.Time          `json:"created_at"`
}

func NewAccountPresenter() AccountPresenter {
	return &accountPresenter{
		included: make([]*Resource, 0),
	}
}

func (p *accountPresenter) Unmarshal(data []byte) (acc model.Account, err error) {
	result := struct {
		Data Account `json:"data"`
	}{}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	return result.Data.toModel(), nil
}

func (p *accountPresenter) Marshal(acc model.Account) ([]byte, error) {
	a := Account{}
	a.fromModel(acc)

	result := struct {
		Data Account `json:"data"`
	}{
		Data: a,
	}

	return json.Marshal(result)
}

func (p *accountPresenter) MarshalList(list []model.Account) ([]byte, error) {
	accounts := make([]Account, 0, len(list))

	for _, acc := range list {
		account := Account{}
		account.fromModel(acc)
		accounts = append(accounts, account)
	}

	result := struct {
		Data []Account `json:"data"`
	}{
		Data: accounts,
	}

	return json.Marshal(result)
}

func (a Account) toModel() model.Account {
	return model.Account{
		ExternalID: a.Attributes.ExternalID,
	}
}

func (a *Account) fromModel(acc model.Account) {
	a.Type.ID = acc.ID
	a.Type.Type = "account"

	a.Attributes.ExternalID = acc.ExternalID
	a.Attributes.Username = acc.Username
	a.Attributes.State = acc.State
	a.Attributes.StateReason = acc.StateReason
	a.Attributes.CreatedAt = acc.CreatedAt
}
