package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/db"
	server "github.com/sefaphlvn/bigbang/pkg/httpserver"
	"github.com/sefaphlvn/bigbang/pkg/log"
	"github.com/sefaphlvn/bigbang/rest/api/auth"
	"github.com/sefaphlvn/bigbang/rest/api/router"
	"github.com/sefaphlvn/bigbang/rest/bridge"
	"github.com/sefaphlvn/bigbang/rest/crud/custom"
	"github.com/sefaphlvn/bigbang/rest/crud/extension"
	"github.com/sefaphlvn/bigbang/rest/crud/scenario"
	"github.com/sefaphlvn/bigbang/rest/crud/xds"
	"github.com/sefaphlvn/bigbang/rest/dependency"
	"github.com/sefaphlvn/bigbang/rest/handlers"
)

// restCmd represents the command for starting the REST API server.
// It initializes the server, sets up routes, and starts listening for incoming HTTP requests.
// Parameters:
// - none
// Returns:
// - *cobra.Command: a Cobra command instance for the REST API server
var restCmd = &cobra.Command{
	Use:   "server-rest",
	Short: "Start Bigbang REST Server",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		appConfig := config.Read(cfgFile)
		logger := log.NewLogger(appConfig)
		appContext := db.NewMongoDB(appConfig, logger, false)
		xdsHandler := xds.NewXDSHandler(appContext)
		extensionHandler := extension.NewExtensionHandler(appContext)
		scenarioHandler := scenario.NewScenarioHandler(appContext)
		customHandler := custom.NewCustomHandler(appContext)
		bridgeHandler := bridge.NewBridgeHandler(appContext)
		userHandler := auth.NewUserHandler(appContext)
		dependencyHandler := dependency.NewDependencyHandler(appContext)
		dependencyHandler.StartCacheCleanup(1 * time.Minute)

		h := handlers.NewHandler(xdsHandler, extensionHandler, customHandler, userHandler, dependencyHandler, bridgeHandler, scenarioHandler)
		r := router.InitRouter(h, logger)
		if err := server.NewHTTPServer(r).Run(appConfig, logger); err != nil {
			logger.Fatalf("Server failed to run: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(restCmd)
}
