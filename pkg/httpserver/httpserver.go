package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/sefaphlvn/bigbang/pkg/config"
)

type Server struct {
	Router *gin.Engine
}

func NewHTTPServer(router *gin.Engine) *Server {
	return &Server{
		Router: router,
	}
}

func (s *Server) Run(config *config.AppConfig, log *logrus.Logger) error {
	server := &http.Server{
		Addr:              ":8099",
		Handler:           s.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Info("Starting http web server [::]:8099")
	<-done
	log.Info("Http web server stop signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	log.Print("Server exited properly")

	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	return err
}
