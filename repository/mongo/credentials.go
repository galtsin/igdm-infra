package mongo

import (
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const credentialsCollectionName = "credentials"

type credentialsRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

type credentials struct {
	ID         primitive.ObjectID `bson:"_id"`
	ExternalID string             `bson:"external_id"`
	Username   string             `bson:"username"`
	Password   string             `bson:"password"`
	Proxy      string             `bson:"proxy"`
}

func CredentialsRepository(db *mongo.Database) domain.CredentialsRepository {
	return &credentialsRepository{
		collection: db.Collection(credentialsCollectionName),
		timeout:    120 * time.Second,
	}
}

func (r *credentialsRepository) Collection() *mongo.Collection {
	return r.collection
}

func (r *credentialsRepository) GetContextTimeout() time.Duration {
	return r.timeout
}

func (r *credentialsRepository) Store(cred model.Credentials) (model.Credentials, error) {
	credentials := credentials{}
	if err := credentials.fromModel(cred); err != nil {
		return model.Credentials{}, err
	}

	var err error
	if cred.ID == "" {
		_, err = insertOne(r, credentials)
	} else {
		_, err = replaceOne(r, &bson.M{"_id": credentials.ID}, credentials)
	}

	if err != nil {
		return model.Credentials{}, err
	}

	return credentials.toModel(), nil
}

func (r *credentialsRepository) Delete(id string) error {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return newErrorInvalidValue(credentialsCollectionName, id, err)
	}

	_, err = deleteOne(r, &bson.M{"_id": bsonID})
	return err
}

func (r *credentialsRepository) WhereExternalID(externalID string) (cred model.Credentials, err error) {
	credentials := credentials{}
	result := findOne(r, &bson.M{"external_id": externalID})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return cred, newErrorNotFound(credentialsCollectionName, externalID)
		}

		return cred, result.Err()
	}

	if err := result.Decode(&credentials); err != nil {
		return cred, err
	}

	return credentials.toModel(), nil
}

func (c *credentials) fromModel(cred model.Credentials) error {
	c.ExternalID = cred.ExternalID
	c.Username = cred.Username
	c.Password = cred.Password
	c.Proxy = cred.Proxy

	if cred.ID == "" {
		c.ID = primitive.NewObjectID()
	} else {
		objectID, err := primitive.ObjectIDFromHex(cred.ID)
		if err != nil {
			return err
		}
		c.ID = objectID
	}

	return nil
}

func (c credentials) toModel() model.Credentials {
	cred := model.Credentials{
		ID:         c.ID.Hex(),
		ExternalID: c.ExternalID,
		Username:   c.Username,
		Password:   c.Password,
		Proxy:      c.Proxy,
	}

	return cred
}
