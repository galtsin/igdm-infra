package domain

import (
	"time"

	"channels-instagram-dm/domain/model"
)

type Repository interface {
	AccountRepository() AccountRepository
	CredentialsRepository() CredentialsRepository
	ConversationRepository() ConversationRepository
	MessageRepository() MessageRepository
	ActivityLogRepository() ActivityLogRepository
}

type AccountRepository interface {
	Store(model.Account) (model.Account, error)
	Delete(id string) error
	All() ([]model.Account, error)
	WhereID(id string) (model.Account, error)
	WhereExternalID(id string) (model.Account, error)
	WhereUsername(username string) (model.Account, error)
	WhereState(state model.AccountState) ([]model.Account, error)
}

type CredentialsRepository interface {
	Store(credentials model.Credentials) (model.Credentials, error)
	Delete(id string) error
	WhereExternalID(externalID string) (model.Credentials, error)
}

type ConversationRepository interface {
	Store(conversation model.Conversation) (model.Conversation, error)
	WhereID(id string) (model.Conversation, error)
	WhereAttributeThreadID(id string) (model.Conversation, error)
}

type ActivityLogRepository interface {
	Store(activityLog model.ActivityLog) (model.ActivityLog, error)
	DeleteWhereAccountIDCreatedAtBefore(accountID string, before time.Time) error
	WhereAccountID(accountID string, limit, offset int) ([]model.ActivityLog, error)
}

type MessageRepository interface {
	Store(message model.Message) (model.Message, error)
	WhereID(id string) (model.Message, error)
	Filter() MessageRepositoryFilter
	InstagramAttributeFilter() MessageRepositoryInstagramAttributeFilter
	WhereChannelsDeliveredFailedRecentAt(filter MessageRepositoryFilter, recentAt time.Duration, limit int) ([]model.Message, error)
	WhereChannelsDeliveredNone(filter MessageRepositoryFilter, limit int) ([]model.Message, error)
	WhereInstagramDeliveredFailedRecentAt(filter MessageRepositoryFilter, recentAt time.Duration, limit int) ([]model.Message, error)
	WhereInstagramDeliveredNone(filter MessageRepositoryFilter, limit int) ([]model.Message, error)
	WhereInstagramAttributeID(id string) (model.Message, error)
	WhereInstagramAttribute(filter MessageRepositoryInstagramAttributeFilter, limit int) ([]model.Message, error)
}

type MessageRepositoryFilter interface {
	WithAccountID(string) MessageRepositoryFilter
	WithSource(model.MessageSource) MessageRepositoryFilter
}

type MessageRepositoryInstagramAttributeFilter interface {
	WithID(string) MessageRepositoryInstagramAttributeFilter
}
