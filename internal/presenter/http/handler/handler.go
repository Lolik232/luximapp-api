package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Lolik232/luximapp-api/internal/domain"
	"github.com/gorilla/mux"
)

type IHandler interface {
	ConfigureRoutes(router *mux.Router)
}

type Handler struct {
}

func (h *Handler) error(w http.ResponseWriter, r *http.Request, err error) {
	httpErr, code := domain.New(err)
	resp := domain.CreateBadResponce(httpErr)
	h.respondJson(w, r, code, resp)
	return
}

func (*Handler) respondJson(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Developed-By","GibbonStudio")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
	return
}
