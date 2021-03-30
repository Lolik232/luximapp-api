package mongo_store

import (
	"context"

	"github.com/Lolik232/luximapp-api/internal/context_keys"
	"github.com/Lolik232/luximapp-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type dblog struct {
	ID     primitive.ObjectID     `bson:"_id"`
	Msg    string                 `bson:"msg"`
	Data   context_keys.InputData `bson:"input_data"`
	Device string                 `bson:"device"`
}

type LogRepository struct {
	database *mongo.Database
	logCol   *mongo.Collection
}

func NewLogRepository(db *mongo.Database) *LogRepository {
	return &LogRepository{
		database: db,
		logCol:   db.Collection(logsCollection),
	}
}

func (l *LogRepository) Create(ctx context.Context, log *domain.Log) error {
	id := primitive.NewObjectID()
	dblog := ToDBLog(log)
	dblog.ID = id
	_, err := l.logCol.InsertOne(ctx, dblog)
	if err != nil {
		return err
	}
	return nil
}

func ToDBLog(log *domain.Log) *dblog {
	return &dblog{
		Msg:    log.Msg,
		Data:   log.Data,
		Device: log.Device,
	}
}
