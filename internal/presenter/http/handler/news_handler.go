package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Lolik232/luximapp-api/internal/domain"
	"github.com/Lolik232/luximapp-api/internal/service"
	"github.com/Lolik232/luximapp-api/pkg/errors"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

type NewsHandler struct {
	Handler
	newsService service.INewsService
}

func NewNewsHandler(svc service.INewsService) *NewsHandler {
	return &NewsHandler{
		Handler:     Handler{},
		newsService: svc,
	}
}

func (n *NewsHandler) ConfigureRoutes(router *mux.Router) {
	news := router.PathPrefix("/news").Subrouter()
	news.HandleFunc("/all", n.all()).Methods(http.MethodGet)
	news.HandleFunc("/fetch", n.fetch()).Methods(http.MethodGet)
	news.HandleFunc("/count", n.count()).Methods(http.MethodGet)
	news.HandleFunc("/new", n.create()).Methods(http.MethodPost)

}
func (n *NewsHandler) count() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := n.newsService.Count(r.Context())
		if err != nil {
			n.error(w, r, err)
			return
		}
		data := map[string]string{
			"count": fmt.Sprintf("%d", count),
		}
		responce := domain.CreateOneOkResponce(data)
		n.respondJson(w, r, 200, responce)
		return
	}
}

func (n *NewsHandler) fetch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// vars := mux.Vars(r)
		c, err := strconv.Atoi(r.FormValue("count"))
		if err != nil {
			err := errors.ErrInvalidArgument.Newf("Invalid count.")
			n.error(w, r, err)
			return
		}
		of, err := strconv.Atoi(r.FormValue("offset"))
		if err != nil {
			err := errors.ErrInvalidArgument.Newf("Invalid offset.")
			n.error(w, r, err)
			return
		}
		if c < 0 {
			err := errors.ErrInvalidArgument.Newf("Count must be positive.")
			n.error(w, r, err)
			return
		}
		if of < 0 {
			err := errors.ErrInvalidArgument.Newf("Offset must be positive.")
			n.error(w, r, err)
			return
		}

		count := uint(c)
		offset := uint(of)
		news, fieldsCount, err := n.newsService.Fetch(r.Context(), offset, count)
		if err != nil {
			n.error(w, r, err)
			return
		}

		response := domain.CreateOkResponce(int(fieldsCount), news)
		n.respondJson(w, r, http.StatusOK, response)
		return
	}
}

func (n *NewsHandler) all() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		news, fieldsCount, err := n.newsService.FindAll(r.Context())
		if err != nil {
			n.error(w, r, err)
			return
		}
		response := domain.CreateOkResponce(int(fieldsCount), news)
		n.respondJson(w, r, http.StatusOK, response)
		return
	}
}

func (n *NewsHandler) create() http.HandlerFunc {
	type CreateInput struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}
	const secretFormKey = "secret_key"

	return func(w http.ResponseWriter, r *http.Request) {

		secretKey := r.FormValue(secretFormKey)
		appSecretKey := os.Getenv("SECRET_KEY")

		if secretKey != appSecretKey {

			n.respondJson(w, r, http.StatusUnauthorized, nil)
			return
		}

		req := &CreateInput{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			err = errors.ErrInvalidArgument.Newf("Invalid input data.")
			n.error(w, r, err)
			return
		}
		news := &domain.News{
			Title: req.Title,
			Text:  req.Text,
		}
		_, err := n.newsService.Create(r.Context(), news)
		if err != nil {
			n.error(w, r, err)
			return
		}
		n.respondJson(w, r, http.StatusCreated, nil)
		return
	}
}
