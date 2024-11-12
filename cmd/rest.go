package cmd

import (
	"time"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/db"
	server "github.com/sefaphlvn/bigbang/pkg/httpserver"
	"github.com/sefaphlvn/bigbang/pkg/log"
	"github.com/sefaphlvn/bigbang/rest/api/auth"
	"github.com/sefaphlvn/bigbang/rest/api/router"
	"github.com/sefaphlvn/bigbang/rest/bridge"
	"github.com/sefaphlvn/bigbang/rest/crud/custom"
	"github.com/sefaphlvn/bigbang/rest/crud/extension"
	"github.com/sefaphlvn/bigbang/rest/crud/xds"
	"github.com/sefaphlvn/bigbang/rest/dependency"
	"github.com/sefaphlvn/bigbang/rest/handlers"

	"github.com/spf13/cobra"
)

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:   "server-rest",
	Short: "Start Bigbang REST Server",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		appConfig := config.Read(cfgFile)
		logger := log.NewLogger(appConfig)
		db := db.NewMongoDB(appConfig, logger)
		xdsHandler := xds.NewXDSHandler(db)
		extensionHandler := extension.NewExtensionHandler(db)
		customHandler := custom.NewCustomHandler(db)
		bridgeHandler := bridge.NewBridgeHandler(db)
		userHandler := auth.NewUserHandler(db)
		dependencyHandler := dependency.NewDependencyHandler(db)
		dependencyHandler.StartCacheCleanup(1 * time.Minute)

		h := handlers.NewHandler(xdsHandler, extensionHandler, customHandler, userHandler, dependencyHandler, bridgeHandler)
		r := router.InitRouter(h, logger)
		if err := server.NewHttpServer(r).Run(appConfig, logger); err != nil {
			logger.Fatalf("Server failed to run: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(restCmd)
}
