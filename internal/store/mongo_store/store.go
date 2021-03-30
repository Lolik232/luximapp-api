package mongo_store

import (
	"github.com/Lolik232/luximapp-api/internal/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	database *mongo.Database
	logRepo  *LogRepository
	newsRepo *NewsRepository
}

func NewStore(db *mongo.Database) store.IStore {
	return &Store{
		database: db,
	}
}

func (st *Store) News() store.INewsRepository {
	if st.newsRepo == nil {
		st.newsRepo = NewNewsRepository(st.database)
	}
	return st.newsRepo
}

func (st *Store) Logs() store.ILogRepository {
	if st.logRepo == nil {
		st.logRepo = NewLogRepository(st.database)
	}
	return st.logRepo
}
