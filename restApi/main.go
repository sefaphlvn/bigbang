package main

import (
	"flag"

	"github.com/sefaphlvn/bigbang/restapi/crud/extension"
	"github.com/sefaphlvn/bigbang/restapi/crud/xds"
	"github.com/sefaphlvn/bigbang/restapi/db"
	"github.com/sefaphlvn/bigbang/restapi/handlers"
	httpserver "github.com/sefaphlvn/bigbang/restapi/http_server"
	"github.com/sefaphlvn/bigbang/restapi/router"
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
