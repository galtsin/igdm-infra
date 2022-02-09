package model

import (
	"fmt"
	"time"
)

const (
	AccountStateActive AccountState = iota + 1
	AccountStateSuspend
)

const (
	AccountStateReasonServiceStopped  = "_STOPPED_BY_SERVICE_" // Остановка сервисом
	AccountStateReasonNoLoggedIn      = "_NO_LOGGED_IN_"
	AccountStateReasonPermanentError  = "_PERMANENT_ERROR_"    // Постоянная ошибка
	AccountStateReasonChallenge       = "_CHALLENGE_REQUIRED_" // Всплыл challenge
	AccountStateReasonMarkedAsDeleted = "_MARKED_AS_DELETED_"  // Пометили на удаление
)

type AccountState int

type AccountStateReason int

type Account struct {
	ID          string
	ExternalID  string
	Username    string
	State       AccountState
	StateReason string
	CreatedAt   time.Time
	InboxSync   InboxSync
}

// LastInboxSyncSnapshot
type InboxSync struct {
	SeqID      int64 // По параметру будут отслеживаться наличие изменений в inbox
	SnapshotAt int64 // По этому параметру будут определяться последние изменения в thread! Если у сообщение timestamp больше
}

func NewAccount(externalID, username string) Account {
	return Account{
		ExternalID: externalID,
		Username:   username,
		State:      AccountStateSuspend,
		CreatedAt:  time.Now(),
		InboxSync:  InboxSync{},
	}
}

func (account Account) Validate() error {
	if account.ExternalID == "" {
		return fmt.Errorf("ExternalID should not be empty")
	}

	return nil
}

func (account *Account) SetInboxSync(inboxSync InboxSync) {
	account.InboxSync = inboxSync
}
