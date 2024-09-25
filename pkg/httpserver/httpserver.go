package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sirupsen/logrus"

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

func (s *Server) Run(config *config.AppConfig, log *logrus.Logger) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.BIGBANG_REST_SERVER_PORT),
		Handler: s.Router,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("Starting http web server [::]:%s", config.BIGBANG_REST_SERVER_PORT)
	<-done
	log.Info("Http web server stop signal recived")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	} else {
		log.Print("Server exited properly")
	}

	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}
