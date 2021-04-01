package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Lolik232/luximapp-api/internal/presenter/http/handler"
	"github.com/Lolik232/luximapp-api/internal/presenter/http/middleware"
	"github.com/Lolik232/luximapp-api/internal/store"
	"github.com/gorilla/mux"
)

type server struct {
	server *http.Server
	router *mux.Router
	store  store.IStore
}

func (s *server) Run() error {
	return s.server.ListenAndServe()
}

func NewServer(st store.IStore, handlers ...handler.IHandler) *server {
	// newsSvc := configureServices(st)

	router := mux.NewRouter()
	configureMiddleware(router)

	for _, h := range handlers {
		h.ConfigureRoutes(router)
	}
	port := os.Getenv("PORT")
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	return &server{
		server: srv,
		router: router,
		store:  st,
	}

}

// func configureServices(store store.IStore) service.INewsService {
// 	newsService := services.NewNewsService(store)
// 	return newsService
// }
func configureMiddleware(router *mux.Router) {
	router.Use(middleware.DeviceMiddleware)
}
