package model

import "channels-instagram-dm/domain/model/instagram"

type Conversation struct {
	ID            string
	AccountID     string
	Attributes    ConversationAttributes
	LastMessageID string
}

type ConversationAttributes struct {
	UserAttributes
	ThreadAttributes
	LastSyncedThreadItemID string // Последнее сохраненное сообщение из треда
}

type ThreadAttributes struct {
	ID               string
	V2ID             string
	LastActivityAt   int64
	Pending          bool
	Archived         bool
	ThreadType       string
	InviterUserID    string
	LastThreadItemID string
}

type UserAttributes struct {
	ID       string
	Username string
	Avatar   string
}

func NewConversation(accountID string) Conversation {
	return Conversation{
		AccountID:  accountID,
		Attributes: ConversationAttributes{},
	}
}

// TODO: Перейти на такой формат
func (c *ConversationAttributes) SetThreadAttributes(thread instagram.Thread) {
	c.ThreadAttributes.ID = thread.ID
	c.ThreadAttributes.V2ID = thread.V2ID
}

func (c *Conversation) SetThreadAttributes(thread instagram.Thread) {
	c.Attributes.ThreadAttributes.ID = thread.ID
	c.Attributes.ThreadAttributes.V2ID = thread.V2ID
	c.Attributes.ThreadAttributes.LastActivityAt = thread.LastActivityAt
	c.Attributes.ThreadAttributes.Pending = thread.Pending
	c.Attributes.ThreadAttributes.Archived = thread.Archived
	c.Attributes.ThreadAttributes.ThreadType = thread.ThreadType
	c.Attributes.ThreadAttributes.InviterUserID = thread.InviterUserID
	c.Attributes.ThreadAttributes.LastThreadItemID = thread.LastPermanentItem.ItemID
}

func (c *Conversation) SetUserAttributes(user instagram.User) {
	c.Attributes.UserAttributes.ID = user.ID
	c.Attributes.UserAttributes.Username = user.Username
	c.Attributes.UserAttributes.Avatar = user.ProfilePicURL
}
