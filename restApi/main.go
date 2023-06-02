package main

import (
	"flag"

	"github.com/sefaphlvn/bigbang/restApi/crud/extension"
	"github.com/sefaphlvn/bigbang/restApi/crud/xds"
	"github.com/sefaphlvn/bigbang/restApi/db"
	"github.com/sefaphlvn/bigbang/restApi/handlers"
	httpserver "github.com/sefaphlvn/bigbang/restApi/http_server"
	"github.com/sefaphlvn/bigbang/restApi/router"
)

func main() {
	flag.Parse()
	db := db.NewMongoDB("mongodb://localhost:27017")
	xdsHandler := xds.NewXDSHandler(db)
	extensionHandler := extension.NewExtensionHandler(db)
	h := handlers.NewHandler(xdsHandler, extensionHandler)
	router := router.InitRouter(h)
	s := httpserver.NewHttpServer(router)
	go s.Run("0.0.0.0:80")

}
