package cmd

import (
	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/httpserver"
	"github.com/sefaphlvn/bigbang/pkg/log"
	"github.com/sefaphlvn/bigbang/restServer/api/auth"
	"github.com/sefaphlvn/bigbang/restServer/api/router"
	"github.com/sefaphlvn/bigbang/restServer/crud/custom"
	"github.com/sefaphlvn/bigbang/restServer/crud/extension"
	"github.com/sefaphlvn/bigbang/restServer/crud/xds"
	"github.com/sefaphlvn/bigbang/restServer/handlers"

	"github.com/spf13/cobra"
)

// restCmd represents the restServer command
var restCmd = &cobra.Command{
	Use:   "server-rest",
	Short: "Start REST Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		var appConfig = config.Read(cfgFile)
		var logger = log.NewLogger(appConfig)
		var db = db.NewMongoDB(appConfig, logger)

		var xdsHandler = xds.NewXDSHandler(db)
		var extensionHandler = extension.NewExtensionHandler(db)
		var customHandler = custom.NewCustomHandler(db)
		var userHandler = auth.NewUserHandler(db)

		h := handlers.NewHandler(xdsHandler, extensionHandler, customHandler, userHandler)
		r := router.InitRouter(h, logger)

		if err := server.NewHttpServer(r).Run(appConfig, logger); err != nil {
			logger.Fatalf("Server failed to run: %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(restCmd)
}
