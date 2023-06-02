package main

import (
	"flag"

	"github.com/sefaphlvn/bigbang/restServer/crud/extension"
	"github.com/sefaphlvn/bigbang/restServer/crud/xds"
	"github.com/sefaphlvn/bigbang/restServer/db"
	"github.com/sefaphlvn/bigbang/restServer/handlers"
	"github.com/sefaphlvn/bigbang/restServer/router"
	httpserver "github.com/sefaphlvn/bigbang/restServer/server"
)

func main() {
	flag.Parse()
	db := db.NewMongoDB("mongodb://localhost:27017")
	xdsHandler := xds.NewXDSHandler(db)
	extensionHandler := extension.NewExtensionHandler(db)
	h := handlers.NewHandler(xdsHandler, extensionHandler)
	router := router.InitRouter(h)
	s := httpserver.NewHttpServer(router)
	s.Run("0.0.0.0:80")
}
