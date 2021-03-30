package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Lolik232/luximapp-api/internal/presenter/http/handler"
	"github.com/Lolik232/luximapp-api/internal/presenter/http/server"
	"github.com/Lolik232/luximapp-api/internal/service/services"
	"github.com/Lolik232/luximapp-api/internal/store/mongo_store"
	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	if err := gotenv.Load(".env"); err != nil {
		log.Println(".env file not found.")
	}
}

func initDb() *mongo.Database {

	databaseURI := os.Getenv("MONGO_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(databaseURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Err in init database. Err message: %s", err.Error())
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Database ping error. %s", err.Error())
	}
	return client.Database("luximapp-api")

}

func main() {
	db := initDb()
	defer db.Client().Disconnect(context.TODO())
	store := mongo_store.NewStore(db)

	newsService := services.NewNewsService(store)
	newsHandler := handler.NewNewsHandler(newsService)

	server := server.NewServer(store, newsHandler)

	log.Println(server.Run())
}
