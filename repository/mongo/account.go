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

const accountCollectionName = "account"

type accountRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

type account struct {
	ID          primitive.ObjectID `bson:"_id"`
	ExternalID  string             `bson:"external_id"`
	Username    string             `bson:"username"`
	State       model.AccountState `bson:"state"`
	StateReason string             `bson:"state_reason"`
	CreatedAt   time.Time          `bson:"created_at"`
	InboxSync   struct {
		SeqID                int64 `bson:"seq_id"`
		PendingRequestsTotal int   `bson:"pending_requests_total"`
		SnapshotAt           int64 `bson:"snapshot_at"`
	} `bson:"inbox_sync"`
}

func AccountRepository(db *mongo.Database) domain.AccountRepository {
	return &accountRepository{
		collection: db.Collection(accountCollectionName),
		timeout:    120 * time.Second,
	}
}

func (r *accountRepository) Collection() *mongo.Collection {
	return r.collection
}

func (r *accountRepository) GetContextTimeout() time.Duration {
	return r.timeout
}

func (r *accountRepository) Store(acc model.Account) (model.Account, error) {
	account := account{}
	if err := account.fromModel(acc); err != nil {
		return model.Account{}, err
	}

	var err error
	if acc.ID == "" {
		_, err = insertOne(r, account)
	} else {
		_, err = replaceOne(r, &bson.M{"_id": account.ID}, account)
	}

	if err != nil {
		return model.Account{}, err
	}

	return account.toModel(), nil
}

func (r *accountRepository) Delete(id string) error {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return newErrorInvalidValue(accountCollectionName, id, err)
	}

	_, err = deleteOne(r, &bson.M{"_id": bsonID})
	return err
}

func (r *accountRepository) All() ([]model.Account, error) {
	result := make([]account, 0)
	err := findAndDecode(r, &bson.M{}, &result)
	if err != nil {
		return nil, err
	}

	accounts := make([]model.Account, 0, len(result))
	for _, account := range result {
		accounts = append(accounts, account.toModel())
	}

	return accounts, nil
}

func (r *accountRepository) WhereID(id string) (acc model.Account, err error) {
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return acc, fmt.Errorf("Invalid value %s with error %s ", id, err)
	}

	account := account{}
	result := findOne(r, &bson.M{"_id": bsonID})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return acc, newErrorNotFound(accountCollectionName, id)
		}

		return acc, result.Err()
	}

	if err := result.Decode(&account); err != nil {
		return acc, err
	}

	return account.toModel(), nil
}

func (r *accountRepository) WhereExternalID(externalID string) (acc model.Account, err error) {
	account := account{}
	result := findOne(r, &bson.M{"external_id": externalID})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return acc, newErrorNotFound(accountCollectionName, externalID)
		}

		return acc, result.Err()
	}

	if err := result.Decode(&account); err != nil {
		return acc, err
	}

	return account.toModel(), nil
}

func (r *accountRepository) WhereState(state model.AccountState) ([]model.Account, error) {
	result := make([]account, 0)
	err := findAndDecode(r, &bson.M{"state": state}, &result)
	if err != nil {
		return nil, err
	}

	accounts := make([]model.Account, 0, len(result))
	for _, account := range result {
		accounts = append(accounts, account.toModel())
	}

	return accounts, nil
}

func (r *accountRepository) WhereUsername(username string) (acc model.Account, err error) {
	account := account{}
	result := findOne(r, &bson.M{"username": username})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return acc, newErrorNotFound(accountCollectionName, username)
		}

		return acc, result.Err()
	}

	if err := result.Decode(&account); err != nil {
		return acc, err
	}

	return account.toModel(), nil
}

func (a *account) fromModel(acc model.Account) error {
	a.ExternalID = acc.ExternalID
	a.Username = acc.Username
	a.State = acc.State
	a.StateReason = acc.StateReason
	a.CreatedAt = acc.CreatedAt
	a.InboxSync.SeqID = acc.InboxSync.SeqID
	a.InboxSync.SnapshotAt = acc.InboxSync.SnapshotAt

	if acc.ID == "" {
		a.ID = primitive.NewObjectID()
	} else {
		objectID, err := primitive.ObjectIDFromHex(acc.ID)
		if err != nil {
			return err
		}
		a.ID = objectID
	}

	return nil
}

func (a account) toModel() model.Account {
	acc := model.Account{
		ID:          a.ID.Hex(),
		ExternalID:  a.ExternalID,
		Username:    a.Username,
		State:       a.State,
		StateReason: a.StateReason,
		CreatedAt:   a.CreatedAt,
		InboxSync: model.InboxSync{
			SeqID:      a.InboxSync.SeqID,
			SnapshotAt: a.InboxSync.SnapshotAt,
		},
	}

	return acc
}
