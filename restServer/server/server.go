package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func (s *Server) Run(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Print("Server Started")

	<-done
	log.Print("Server Stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	} else {
		log.Print("Server Exited Properly")
	}

	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}
