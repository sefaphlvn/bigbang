package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func NewHttpServer(router *gin.Engine) *Server {
	return &Server{
		Router: router,
	}
}

func (s *Server) Run(addr string) {
	server := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
