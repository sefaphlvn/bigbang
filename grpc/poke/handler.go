package poke

import (
	"encoding/json"
	"fmt"
	"github.com/sefaphlvn/bigbang/grpc/server"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/sefaphlvn/bigbang/pkg/db"
)

type Handler struct {
	Ctx    *server.Context
	DB     *db.MongoDB
	Logger *logrus.Logger
	Func   server.Func
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ping":
		h.handlePing(w, r)
	case "/poke":
		h.handlePoke(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) handlePing(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info(fmt.Fprint(w, "OK"))
}

func (h *Handler) handlePoke(w http.ResponseWriter, r *http.Request) {
	allResources, err := h.Func.GetAllResourcesFromListener("newListener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err)
		return
	}

	err = h.Ctx.SetSnapshot(allResources, h.Logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(allResources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err)
		return
	}
}
