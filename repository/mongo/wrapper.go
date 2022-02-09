package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Collection() *mongo.Collection
	GetContextTimeout() time.Duration
}

func findOne(r Repository, query interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetContextTimeout())
	defer cancel()
	return r.Collection().FindOne(ctx, query, opts...)
}

func insertOne(r Repository, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetContextTimeout())
	defer cancel()
	return r.Collection().InsertOne(ctx, document, opts...)
}

func insertMany(r Repository, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetContextTimeout())
	defer cancel()
	return r.Collection().InsertMany(ctx, documents, opts...)
}

func replaceOne(r Repository, query interface{}, document interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetContextTimeout())
	defer cancel()
	return r.Collection().ReplaceOne(ctx, query, document, opts...)
}

func updateMany(r Repository, query interface{}, data interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetContextTimeout())
	defer cancel()
	return r.Collection().UpdateMany(ctx, query, data, opts...)
}

func deleteOne(r Repository, query interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetContextTimeout())
	defer cancel()
	return r.Collection().DeleteOne(ctx, query, opts...)
}

func deleteMany(r Repository, query interface{}, opts ...*options.DeleteOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetContextTimeout())
	defer cancel()
	_, err := r.Collection().DeleteMany(ctx, query, opts...)

	return err
}

func countDocuments(r Repository, query interface{}, opts ...*options.CountOptions) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetContextTimeout())
	defer cancel()
	return r.Collection().CountDocuments(ctx, query, opts...)
}

func findAndDecode(r Repository, query interface{}, results interface{}, opts ...*options.FindOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetContextTimeout())
	defer cancel()

	cur, err := r.Collection().Find(ctx, query, opts...)
	if err != nil {
		return err
	}

	return cur.All(ctx, results)
}
