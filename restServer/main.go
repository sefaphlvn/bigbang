package main

import (
	"log"

	"github.com/sefaphlvn/bigbang/restServer/auth"
	"github.com/sefaphlvn/bigbang/restServer/crud/custom"
	"github.com/sefaphlvn/bigbang/restServer/crud/extension"
	"github.com/sefaphlvn/bigbang/restServer/crud/xds"
	"github.com/sefaphlvn/bigbang/restServer/db"
	"github.com/sefaphlvn/bigbang/restServer/handlers"
	"github.com/sefaphlvn/bigbang/restServer/router"
	httpserver "github.com/sefaphlvn/bigbang/restServer/server"
)

func main() {
	db, err := db.NewMongoDB("mongodb://localhost:27017")

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	xdsHandler := xds.NewXDSHandler(db)
	extensionHandler := extension.NewExtensionHandler(db)
	customHandler := custom.NewCustomHandler(db)
	userHandler := auth.NewUserHandler(db)
	h := handlers.NewHandler(xdsHandler, extensionHandler, customHandler, userHandler)

	// Router initialization
	router := router.InitRouter(h)

	s := httpserver.NewHttpServer(router)
	if err := s.Run("0.0.0.0:80"); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}

}
