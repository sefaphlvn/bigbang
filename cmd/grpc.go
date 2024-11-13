package cmd

import (
	"context"

	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/spf13/cobra"

	"github.com/sefaphlvn/bigbang/grpc/poke"
	grpcserver "github.com/sefaphlvn/bigbang/grpc/server"
	"github.com/sefaphlvn/bigbang/grpc/server/bridge"
	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/log"
)

var (
	port   uint
	nodeID string
)

// grpcCmd represents the command for starting the gRPC server.
// It initializes the server, sets up the necessary services, and starts listening for incoming gRPC requests.
// Parameters:
// - none
// Returns:
// - *cobra.Command: a Cobra command instance for the gRPC server
var grpcCmd = &cobra.Command{
	Use:   "server-grpc",
	Short: "Start Bigbang GRPC Server",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		appConfig := config.Read(cfgFile)
		logger := log.NewLogger(appConfig)
		db := db.NewMongoDB(appConfig, logger)
		ctxCache := snapshot.GetContext(logger)

		pokeServer := poke.NewPokeServer(ctxCache, db, logger, appConfig)
		go pokeServer.Run()
		errorContext := bridge.NewErrorContext(10)

		callbacks := grpcserver.NewCallbacks(logger, errorContext)
		srv := server.NewServer(context.Background(), ctxCache.Cache.Cache, callbacks)
		grpcServer := grpcserver.NewServer(srv, port, logger, ctxCache)

		grpcServer.Run(db, errorContext)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.PersistentFlags().UintVar(&port, "port", 18000, "xDS management server port")
	grpcCmd.PersistentFlags().StringVar(&nodeID, "nodeID", "test", "Node ID")
}
