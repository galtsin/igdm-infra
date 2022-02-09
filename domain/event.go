package domain

import (
	"channels-instagram-dm/domain/model"
)

// Есть 2 типа событий:
// 1. Уведомление об изменении в системе
// 2. Требование на необходимость что-то сделать

type EventBus interface {
	PublishAccountCreated(EventAccountCreated)
	PublishAccountResumed(EventAccountResumed)
	PublishAccountSuspended(EventAccountSuspended)
	PublishAccountDeleted(EventAccountDeleted)
	PublishAccountLogout(EventAccountLogout)
	PublishSuspendAccount(EventSuspendAccount)
	PublishInboxHasChanges(EventInboxHasChanges)
	PublishLoginAccount(EventLoginAccount)
	SubscribeOnAccountCreated(EventFilter) chan EventAccountCreated
	SubscribeOnAccountResumed(EventFilter) chan EventAccountResumed
	SubscribeOnAccountSuspended(EventFilter) chan EventAccountSuspended
	SubscribeOnAccountDeleted(EventFilter) chan EventAccountDeleted
	SubscribeOnAccountLogout(EventFilter) chan EventAccountLogout
	SubscribeOnSuspendAccount(EventFilter) chan EventSuspendAccount
	SubscribeOnInboxHasChanges(EventFilter) chan EventInboxHasChanges
	SubscribeOnLoginAccount(EventFilter) chan EventLoginAccount
}

type EventFilter func(event interface{}) bool

type EventAccountCreated struct {
	Account model.Account
}

type EventAccountResumed struct {
	Account model.Account
}

type EventAccountSuspended struct {
	Account model.Account
}

type EventAccountDeleted struct {
	Account model.Account
}

type EventAccountLogout struct {
	Account model.Account
}

type EventSuspendAccount struct {
	Reason  string
	Account model.Account
}

type EventInboxHasChanges struct {
	Account model.Account
}

type EventLoginAccount struct {
	Account model.Account
}
