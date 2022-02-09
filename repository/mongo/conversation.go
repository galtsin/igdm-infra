package mongo

import (
	"fmt"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const conversationCollectionName = "conversation"

type conversationRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

type conversation struct {
	ID            primitive.ObjectID     `bson:"_id"`
	AccountID     string                 `bson:"account_id"`
	Attributes    ConversationAttributes `bson:"attributes"`
	LastMessageID string                 `bson:"last_message_id"`
}

type ConversationAttributes struct {
	UserAttributes         `bson:"user"`
	ThreadAttributes       `bson:"thread"`
	LastSyncedThreadItemID string `bson:"last_synced_thread_item_id"`
}

type UserAttributes struct {
	UserID   string `bson:"id"`
	Username string `bson:"username"`
	Avatar   string `bson:"avatar"`
}

type ThreadAttributes struct {
	ID               string `bson:"id"`
	V2ID             string `bson:"v2_id"`
	LastActivityAt   int64  `bson:"last_activity_at"`
	Pending          bool   `bson:"pending"`
	Archived         bool   `bson:"archived"`
	ThreadType       string `bson:"thread_type"`
	InviterUserID    string `bson:"inviter"`
	LastThreadItemID string `bson:"last_thread_item_id"`
}

func ConversationRepository(db *mongo.Database) domain.ConversationRepository {
	return &conversationRepository{
		collection: db.Collection(conversationCollectionName),
		timeout:    120 * time.Second,
	}
}

func (r *conversationRepository) Collection() *mongo.Collection {
	return r.collection
}

func (r *conversationRepository) GetContextTimeout() time.Duration {
	return r.timeout
}

func (r *conversationRepository) Store(conv model.Conversation) (model.Conversation, error) {
	conversation := conversation{}
	if err := conversation.fromModel(conv); err != nil {
		return model.Conversation{}, err
	}

	var err error
	if conv.ID == "" {
		_, err = insertOne(r, conversation)
	} else {
		_, err = replaceOne(r, &bson.M{"_id": conversation.ID}, conversation)
	}

	if err != nil {
		return model.Conversation{}, err
	}

	return conversation.toModel(), nil
}

func (r *conversationRepository) Delete(id string) error {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return newErrorInvalidValue(conversationCollectionName, id, err)
	}

	_, err = deleteOne(r, &bson.M{"_id": bsonID})
	return err
}

func (r *conversationRepository) WhereID(id string) (conv model.Conversation, err error) {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return conv, fmt.Errorf("Invalid value %s with error %s ", id, err)
	}

	conversation := conversation{}
	result := findOne(r, &bson.M{"_id": bsonID})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return conv, newErrorNotFound(conversationCollectionName, id)
		}

		return conv, result.Err()
	}

	if err := result.Decode(&conversation); err != nil {
		return conv, err
	}

	return conversation.toModel(), nil
}

func (r *conversationRepository) WhereAttributeThreadID(id string) (conv model.Conversation, err error) {
	conversation := conversation{}
	result := findOne(r, &bson.M{"attributes.thread.id": id})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return conv, newErrorNotFound(conversationCollectionName, id)
		}

		return conv, result.Err()
	}

	if err := result.Decode(&conversation); err != nil {
		return conv, err
	}

	return conversation.toModel(), nil
}

func (c *conversation) fromModel(conv model.Conversation) error {
	if conv.ID == "" {
		c.ID = primitive.NewObjectID()
	} else {
		objectID, err := primitive.ObjectIDFromHex(conv.ID)
		if err != nil {
			return err
		}
		c.ID = objectID
	}

	c.AccountID = conv.AccountID
	c.LastMessageID = conv.LastMessageID
	c.Attributes = ConversationAttributes{
		UserAttributes: UserAttributes{
			UserID:   conv.Attributes.UserAttributes.ID,
			Username: conv.Attributes.UserAttributes.Username,
			Avatar:   conv.Attributes.UserAttributes.Avatar,
		},
		ThreadAttributes: ThreadAttributes{
			ID:               conv.Attributes.ThreadAttributes.ID,
			V2ID:             conv.Attributes.ThreadAttributes.V2ID,
			LastActivityAt:   conv.Attributes.ThreadAttributes.LastActivityAt,
			Pending:          conv.Attributes.ThreadAttributes.Pending,
			Archived:         conv.Attributes.ThreadAttributes.Archived,
			ThreadType:       conv.Attributes.ThreadAttributes.ThreadType,
			InviterUserID:    conv.Attributes.ThreadAttributes.InviterUserID,
			LastThreadItemID: conv.Attributes.ThreadAttributes.LastThreadItemID,
		},
		LastSyncedThreadItemID: "",
	}

	return nil
}

func (c conversation) toModel() model.Conversation {
	conv := model.Conversation{
		ID:            c.ID.Hex(),
		AccountID:     c.AccountID,
		LastMessageID: c.LastMessageID,
		Attributes: model.ConversationAttributes{
			UserAttributes: model.UserAttributes{
				ID:       c.Attributes.UserAttributes.UserID,
				Username: c.Attributes.UserAttributes.Username,
				Avatar:   c.Attributes.UserAttributes.Avatar,
			},
			ThreadAttributes: model.ThreadAttributes{
				ID:               c.Attributes.ThreadAttributes.ID,
				V2ID:             c.Attributes.ThreadAttributes.V2ID,
				LastActivityAt:   c.Attributes.ThreadAttributes.LastActivityAt,
				Pending:          c.Attributes.ThreadAttributes.Pending,
				Archived:         c.Attributes.ThreadAttributes.Archived,
				ThreadType:       c.Attributes.ThreadAttributes.ThreadType,
				InviterUserID:    c.Attributes.ThreadAttributes.InviterUserID,
				LastThreadItemID: c.Attributes.ThreadAttributes.LastThreadItemID,
			},
		},
	}

	return conv
}
