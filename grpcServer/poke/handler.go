package poke

import (
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
	h.Func.GetConfigurationFromListener("sad")
	fmt.Fprint(w, "pong")
}

func (h *Handler) handlePoke(w http.ResponseWriter, r *http.Request) {
	// Burada h.Ctx örneğinizi kullanabilirsiniz.
	fmt.Fprint(w, "poke")
}
