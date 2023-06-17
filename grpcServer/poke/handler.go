package poke

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sefaphlvn/bigbang/grpcServer/db"
	grpcserver "github.com/sefaphlvn/bigbang/grpcServer/server"
)

type Handler struct {
	Ctx  *grpcserver.Context
	DB   *db.MongoDB
	L    *grpcserver.Logger
	Func grpcserver.Func
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
	fmt.Fprint(w, "OK")
}

func (h *Handler) handlePoke(w http.ResponseWriter, r *http.Request) {
	asd, err := h.Func.GetAllResourcesFromListener("sefa")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Ctx.SetSnapshot(asd, *h.L)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(asd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
