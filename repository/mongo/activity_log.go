package mongo

import (
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const activityLogCollectionName = "activity_log"

type activityLogRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

type activityLog struct {
	ID        primitive.ObjectID `bson:"_id"`
	AccountID string             `bson:"account_id"`
	Log       string             `bson:"log"`
	CreatedAt time.Time          `bson:"created_at"`
}

func ActivityLogRepository(db *mongo.Database) domain.ActivityLogRepository {
	return &activityLogRepository{
		collection: db.Collection(activityLogCollectionName),
		timeout:    120 * time.Second,
	}
}

func (r *activityLogRepository) Collection() *mongo.Collection {
	return r.collection
}

func (r *activityLogRepository) GetContextTimeout() time.Duration {
	return r.timeout
}

func (r *activityLogRepository) Store(al model.ActivityLog) (model.ActivityLog, error) {
	dbModel := activityLog{}
	if err := dbModel.fromModel(al); err != nil {
		return model.ActivityLog{}, err
	}

	_, err := insertOne(r, dbModel)
	if err != nil {
		return model.ActivityLog{}, err
	}

	return dbModel.toModel(), nil
}

func (r *activityLogRepository) WhereAccountID(accountID string, limit, offset int) ([]model.ActivityLog, error) {
	var dbResult []activityLog

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"created_at": -1})

	err := findAndDecode(r, bson.M{"account_id": accountID}, &dbResult, findOptions)
	if err != nil {
		return nil, err
	}

	result := make([]model.ActivityLog, 0, len(dbResult))
	for _, r := range dbResult {
		result = append(result, r.toModel())
	}

	return result, nil
}

func (r *activityLogRepository) DeleteWhereAccountIDCreatedAtBefore(accountID string, before time.Time) error {
	return deleteMany(r, &bson.M{"account_id": accountID, "created_at": &bson.M{"$lt": before}})
}

func (c *activityLog) fromModel(activityLog model.ActivityLog) error {
	c.AccountID = activityLog.AccountID
	c.Log = activityLog.Log
	c.CreatedAt = activityLog.CreatedAt
	c.ID = primitive.NewObjectID()

	return nil
}

func (c activityLog) toModel() model.ActivityLog {
	return model.ActivityLog{
		AccountID: c.AccountID,
		Log:       c.Log,
		CreatedAt: c.CreatedAt,
	}
}
