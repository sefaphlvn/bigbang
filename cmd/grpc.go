package cmd

import (
	"context"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/sefaphlvn/bigbang/grpc/poke"
	grpcserver "github.com/sefaphlvn/bigbang/grpc/server"
	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/log"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	port   uint
	nodeID string
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "server-grpc",
	Short: "Start GRPC Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		var appConfig = config.Read(cfgFile)
		var logger = log.NewLogger(appConfig)
		var db = db.NewMongoDB(appConfig, logger)

		// Create cache
		ctxCache := grpcserver.GetContext(logger)
		grpcServerHandler := &grpcserver.Handler{Ctx: ctxCache, DB: db, Logger: logger}
		pokeHandler := &poke.Handler{Ctx: ctxCache, DB: db, Logger: logger, Func: grpcServerHandler}

		// Start http server
		go func() {
			err := http.ListenAndServe(":8080", pokeHandler)
			if err != nil {
				logger.Fatalf("failed to start HTTP server: %v", err)
			}
		}()

		// Set initial snapshots
		grpcServerHandler.InitialSnapshots()
		logger.Infof("All snapshots are loaded")

		// Start grpc server
		ctx := context.Background()
		cb := &grpcserver.Callbacks{Logger: logger}
		srv := server.NewServer(ctx, ctxCache.Cash.Cache, cb)
		grpcserver.RunServer(srv, port)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.PersistentFlags().UintVar(&port, "port", 18000, "xDS management server port")
	grpcCmd.PersistentFlags().StringVar(&nodeID, "nodeID", "test", "Node ID")
}
