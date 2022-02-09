package mongo

import (
	"errors"
	"fmt"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const messageCollectionName = "message"

type messageRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

type message struct {
	ID             primitive.ObjectID  `bson:"_id"`
	AccountID      string              `bson:"account_id"`
	ConversationID string              `bson:"conversation_id"`
	Source         model.MessageSource `bson:"source"`
	Type           model.MessageType   `bson:"type"`
	Payload        MessagePayload      `bson:"payload"`
	Attributes     MessageAttributes   `bson:"attributes"`
	Delivered      MessageDelivered    `bson:"delivered"`
	CreatedAt      time.Time           `bson:"created_at"`
}

type MessageAttributes struct {
	InstagramAttributes `bson:"instagram,omitempty"`
	ChannelsAttributes  `bson:"channels,omitempty"`
}

type InstagramAttributes struct {
	ID        string `bson:"id"`
	UserID    string `bson:"user_id"`
	Timestamp int64  `bson:"timestamp"`
}

type ChannelsAttributes struct {
	ID string `bson:"id"`
}

type MessageDelivered struct {
	Status    model.MessageDeliveryStatus `bson:"status"`
	AttemptAt time.Time                   `bson:"attempt_at"`
}

type MessagePayload struct {
	ID              string `bson:"id,omitempty"`
	Text            string `bson:"text,omitempty"`
	Like            string `bson:"like,omitempty"`
	Width           int    `bson:"width,omitempty"`
	Height          int    `bson:"height,omitempty"`
	Url             string `bson:"url,omitempty"`
	Title           string `bson:"title,omitempty"`
	Summary         string `bson:"summary,omitempty"`
	ImagePreviewUrl string `bson:"image_preview_url,omitempty"`
}

func MessageRepository(db *mongo.Database) domain.MessageRepository {
	return &messageRepository{
		collection: db.Collection(messageCollectionName),
		timeout:    120 * time.Second,
	}
}

func (r *messageRepository) Collection() *mongo.Collection {
	return r.collection
}

func (r *messageRepository) GetContextTimeout() time.Duration {
	return r.timeout
}

func (r *messageRepository) Store(msg model.Message) (model.Message, error) {
	var message message
	if err := message.fromModel(msg); err != nil {
		return model.Message{}, err
	}

	var err error
	if msg.ID == "" {
		_, err = insertOne(r, message)
	} else {
		_, err = replaceOne(r, bson.M{"_id": message.ID}, message)
	}

	if err != nil {
		return model.Message{}, err
	}

	return message.toModel(), nil
}

func (r *messageRepository) Delete(id string) error {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return newErrorInvalidValue(messageCollectionName, id, err)
	}

	_, err = deleteOne(r, bson.M{"_id": bsonID})
	return err
}

func (r *messageRepository) WhereID(id string) (model.Message, error) {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Message{}, fmt.Errorf("Invalid value %s with error %s ", id, err)
	}

	var dbResult message

	result := findOne(r, bson.M{"_id": bsonID})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return model.Message{}, newErrorNotFound(messageCollectionName, id)
		}

		return model.Message{}, result.Err()
	}

	if err := result.Decode(&dbResult); err != nil {
		return model.Message{}, err
	}

	return dbResult.toModel(), nil
}

func (r *messageRepository) WhereChannelsDeliveredNone(filter domain.MessageRepositoryFilter, limit int) ([]model.Message, error) {
	var dbResult []message

	f, ok := filter.(*MessageRepositoryFilter)
	if !ok {
		return nil, fmt.Errorf("Filter has wrong type")
	}

	f.WithSource(model.MessageSourceChannels)

	query := f.toMap()
	query["delivered.status"] = model.MessageDeliveryStatusNone

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": 1})

	err := findAndDecode(r, query, &dbResult, findOptions)
	if err != nil {
		return nil, err
	}

	result := make([]model.Message, 0, len(dbResult))
	for _, r := range dbResult {
		result = append(result, r.toModel())
	}

	return result, nil
}

func (r *messageRepository) WhereInstagramDeliveredNone(filter domain.MessageRepositoryFilter, limit int) ([]model.Message, error) {
	var dbResult []message

	f, ok := filter.(*MessageRepositoryFilter)
	if !ok {
		return nil, fmt.Errorf("Filter has wrong type")
	}

	f.WithSource(model.MessageSourceInstagram)

	query := f.toMap()
	query["delivered.status"] = model.MessageDeliveryStatusNone

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"attributes.instagram.timestamp": 1})

	err := findAndDecode(r, query, &dbResult, findOptions)
	if err != nil {
		return nil, err
	}

	result := make([]model.Message, 0, len(dbResult))
	for _, r := range dbResult {
		result = append(result, r.toModel())
	}

	return result, nil
}

func (r *messageRepository) WhereChannelsDeliveredFailedRecentAt(filter domain.MessageRepositoryFilter, recentAt time.Duration, limit int) ([]model.Message, error) {
	var dbResult []message

	f, ok := filter.(*MessageRepositoryFilter)
	if !ok {
		return nil, fmt.Errorf("Filter has wrong type")
	}

	f.WithSource(model.MessageSourceChannels)

	query := f.toMap()
	query["delivered.status"] = model.MessageDeliveryStatusFailed
	query["created_at"] = bson.M{"$gte": time.Now().Add(-1 * recentAt)}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": 1})

	err := findAndDecode(r, query, &dbResult, findOptions)
	if err != nil {
		return nil, err
	}

	result := make([]model.Message, 0, len(dbResult))
	for _, r := range dbResult {
		result = append(result, r.toModel())
	}

	return result, nil
}

func (r *messageRepository) WhereInstagramDeliveredFailedRecentAt(filter domain.MessageRepositoryFilter, recentAt time.Duration, limit int) ([]model.Message, error) {
	var dbResult []message

	f, ok := filter.(*MessageRepositoryFilter)
	if !ok {
		return nil, fmt.Errorf("Filter has wrong type")
	}

	f.WithSource(model.MessageSourceInstagram)

	query := f.toMap()
	query["delivered.status"] = model.MessageDeliveryStatusFailed
	query["created_at"] = bson.M{"$gte": time.Now().Add(-1 * recentAt)}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"attributes.instagram.timestamp": 1})

	err := findAndDecode(r, query, &dbResult, findOptions)
	if err != nil {
		return nil, err
	}

	result := make([]model.Message, 0, len(dbResult))
	for _, r := range dbResult {
		result = append(result, r.toModel())
	}

	return result, nil
}

func (r *messageRepository) WhereInstagramAttributeID(id string) (msg model.Message, err error) {
	var dbResult message

	result := findOne(r, bson.M{"attributes.instagram.id": id})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return msg, newErrorNotFound(messageCollectionName, id)
		}

		return msg, result.Err()
	}

	if err := result.Decode(&dbResult); err != nil {
		return msg, err
	}

	return dbResult.toModel(), nil
}

func (r *messageRepository) WhereInstagramAttribute(filter domain.MessageRepositoryInstagramAttributeFilter, limit int) ([]model.Message, error) {
	var dbResult []message

	f, ok := filter.(*MessageRepositoryInstagramAttributeFilter)
	if !ok {
		return nil, fmt.Errorf("Filter has wrong type")
	}

	query := f.toMap()

	findOptions := options.Find().SetLimit(int64(limit))

	err := findAndDecode(r, query, &dbResult, findOptions)
	if err != nil {
		return nil, err
	}

	result := make([]model.Message, 0, len(dbResult))
	for _, r := range dbResult {
		result = append(result, r.toModel())
	}

	return result, nil
}

type MessageRepositoryFilter struct {
	accountID []string
	source    model.MessageSource
}

func (r *messageRepository) Filter() domain.MessageRepositoryFilter {
	return &MessageRepositoryFilter{
		accountID: make([]string, 0),
	}
}

func (f *MessageRepositoryFilter) WithAccountID(id string) domain.MessageRepositoryFilter {
	f.accountID = append(f.accountID, id)
	return f
}

func (f *MessageRepositoryFilter) WithSource(source model.MessageSource) domain.MessageRepositoryFilter {
	f.source = source
	return f
}

func (f *MessageRepositoryFilter) toMap() bson.M {
	filter := bson.M{}

	if len(f.accountID) != 0 {
		filter["account_id"] = bson.M{"$in": f.accountID}
	}

	if f.source != "" {
		filter["source"] = f.source
	}

	return filter
}

type MessageRepositoryInstagramAttributeFilter struct {
	ID []string
}

func (r *messageRepository) InstagramAttributeFilter() domain.MessageRepositoryInstagramAttributeFilter {
	return &MessageRepositoryInstagramAttributeFilter{}
}

func (f *MessageRepositoryInstagramAttributeFilter) WithID(id string) domain.MessageRepositoryInstagramAttributeFilter {
	f.ID = append(f.ID, id)
	return f
}

func (f *MessageRepositoryInstagramAttributeFilter) toMap() bson.M {
	filter := bson.M{}

	if len(f.ID) != 0 {
		filter["_id"] = bson.M{"$in": f.ID}
	}

	return filter
}

func (m message) toModel() model.Message {
	msg := model.Message{
		ID:             m.ID.Hex(),
		AccountID:      m.AccountID,
		ConversationID: m.ConversationID,
		Source:         m.Source,
		Type:           m.Type,
		CreatedAt:      m.CreatedAt,
		Delivered: model.MessageDelivered{
			Status:    m.Delivered.Status,
			AttemptAt: m.Delivered.AttemptAt,
		},
		Attributes: model.MessageAttributes{
			InstagramAttributes: model.InstagramAttributes{
				ID:        m.Attributes.InstagramAttributes.ID,
				UserID:    m.Attributes.InstagramAttributes.UserID,
				Timestamp: m.Attributes.InstagramAttributes.Timestamp,
			},
			ChannelsAttributes: model.ChannelsAttributes{
				ID: m.Attributes.ChannelsAttributes.ID,
			},
		},
	}

	switch m.Type {
	case model.MessageTypeText:
		msg.Payload = model.MessageText{
			Text: m.Payload.Text,
		}
	case model.MessageTypeLike:
		msg.Payload = model.MessageLike{
			Like: m.Payload.Like,
		}
	case model.MessageTypeLink:
		msg.Payload = model.MessageLink{
			Url:             m.Payload.Url,
			Title:           m.Payload.Title,
			Summary:         m.Payload.Summary,
			ImagePreviewUrl: m.Payload.ImagePreviewUrl,
		}
	case model.MessageTypeActionLog:
		msg.Payload = model.MessageActionLog{
			Text: m.Payload.Text,
		}
	case model.MessageTypeMediaImage:
		msg.Payload = model.MessageMediaImage{
			MessageMedia: model.MessageMedia{
				ID:  m.Payload.ID,
				Url: m.Payload.Url,
			},
			Width:  m.Payload.Width,
			Height: m.Payload.Height,
		}
	case model.MessageTypeMediaVideo:
		msg.Payload = model.MessageMediaVideo{
			MessageMedia: model.MessageMedia{
				ID:  m.Payload.ID,
				Url: m.Payload.Url,
			},
			Width:  m.Payload.Width,
			Height: m.Payload.Height,
		}
	case model.MessageTypeMediaVisualImage:
		msg.Payload = model.MessageMediaVisualImage{
			MessageMedia: model.MessageMedia{
				ID:  m.Payload.ID,
				Url: m.Payload.Url,
			},
			Width:  m.Payload.Width,
			Height: m.Payload.Height,
		}
	case model.MessageTypeMediaVisualVideo:
		msg.Payload = model.MessageMediaVisualVideo{
			MessageMedia: model.MessageMedia{
				ID:  m.Payload.ID,
				Url: m.Payload.Url,
			},
			Width:  m.Payload.Width,
			Height: m.Payload.Height,
		}
	case model.MessageTypeMediaAnimated:
		msg.Payload = model.MessageMediaAnimated{
			MessageMedia: model.MessageMedia{
				ID:  m.Payload.ID,
				Url: m.Payload.Url,
			},
			Width:  m.Payload.Width,
			Height: m.Payload.Height,
		}
	case model.MessageTypeMediaVoice:
		msg.Payload = model.MessageMediaVoice{
			MessageMedia: model.MessageMedia{
				ID:  m.Payload.ID,
				Url: m.Payload.Url,
			},
		}
	case model.MessageTypeUndefined:
		msg.Payload = model.MessageText{
			Text: m.Payload.Text,
		}
	}

	return msg
}

func (m *message) fromModel(msg model.Message) error {
	if msg.ID == "" {
		m.ID = primitive.NewObjectID()
	} else {
		objectID, err := primitive.ObjectIDFromHex(msg.ID)
		if err != nil {
			return err
		}
		m.ID = objectID
	}

	m.AccountID = msg.AccountID
	m.ConversationID = msg.ConversationID
	m.Source = msg.Source
	m.Type = msg.Type
	m.CreatedAt = msg.CreatedAt
	m.Delivered = MessageDelivered{
		Status:    msg.Delivered.Status,
		AttemptAt: msg.Delivered.AttemptAt,
	}
	m.Attributes = MessageAttributes{
		InstagramAttributes: InstagramAttributes{
			ID:        msg.Attributes.InstagramAttributes.ID,
			UserID:    msg.Attributes.InstagramAttributes.UserID,
			Timestamp: msg.Attributes.InstagramAttributes.Timestamp,
		},
		ChannelsAttributes: ChannelsAttributes{
			ID: msg.Attributes.ChannelsAttributes.ID,
		},
	}

	if msg.Payload == nil {
		return errors.New("Payload is empty")
	}

	switch payload := msg.Payload.(type) {
	case model.MessageText:
		m.Payload.Text = payload.Text
	case model.MessageLike:
		m.Payload.Like = payload.Like
	case model.MessageLink:
		m.Payload.Url = payload.Url
		m.Payload.Title = payload.Title
		m.Payload.Summary = payload.Summary
		m.Payload.ImagePreviewUrl = payload.ImagePreviewUrl
	case model.MessageActionLog:
		m.Payload.Text = payload.Text
	case model.MessageMediaImage:
		m.Payload.ID = payload.ID
		m.Payload.Url = payload.Url
		m.Payload.Width = payload.Width
		m.Payload.Height = payload.Height
	case model.MessageMediaVideo:
		m.Payload.ID = payload.ID
		m.Payload.Url = payload.Url
		m.Payload.Width = payload.Width
		m.Payload.Height = payload.Height
	case model.MessageMediaVisualImage:
		m.Payload.ID = payload.ID
		m.Payload.Url = payload.Url
		m.Payload.Width = payload.Width
		m.Payload.Height = payload.Height
	case model.MessageMediaVisualVideo:
		m.Payload.ID = payload.ID
		m.Payload.Url = payload.Url
		m.Payload.Width = payload.Width
		m.Payload.Height = payload.Height
	case model.MessageMediaAnimated:
		m.Payload.ID = payload.ID
		m.Payload.Url = payload.Url
		m.Payload.Width = payload.Width
		m.Payload.Height = payload.Height
	case model.MessageMediaVoice:
		m.Payload.ID = payload.ID
		m.Payload.Url = payload.Url
	case model.MessageUndefined:
		m.Payload.Text = payload.Text
	}

	return nil
}
