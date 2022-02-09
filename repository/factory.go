package repository

import (
	"context"
	"time"

	"channels-instagram-dm/domain"
	mongoRepository "channels-instagram-dm/repository/mongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const defaultTimeout = 10 * time.Second

type factory struct {
	logger domain.Logger
	db     *mongo.Database
}

func Factory(ctx context.Context, logger domain.Logger, url string, dbName string) (domain.Repository, error) {
	ctxOnQuery, cancel := context.WithTimeout(ctx, defaultTimeout)

	go func() {
		<-ctx.Done()
		cancel()
	}()

	client, err := mongo.Connect(ctxOnQuery, options.Client().ApplyURI(url).SetReadPreference(readpref.Primary()))
	if err != nil {
		return nil, err
	}

	factory := &factory{
		logger: logger,
		db:     client.Database(dbName),
	}

	return factory, nil
}

func (f *factory) AccountRepository() domain.AccountRepository {
	return mongoRepository.AccountRepository(f.db)
}

func (f *factory) CredentialsRepository() domain.CredentialsRepository {
	return mongoRepository.CredentialsRepository(f.db)
}

func (f *factory) ConversationRepository() domain.ConversationRepository {
	return mongoRepository.ConversationRepository(f.db)
}

func (f *factory) MessageRepository() domain.MessageRepository {
	return mongoRepository.MessageRepository(f.db)
}

func (f *factory) ActivityLogRepository() domain.ActivityLogRepository {
	return mongoRepository.ActivityLogRepository(f.db)
}
