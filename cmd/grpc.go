package cmd

import (
	"context"

	"github.com/sefaphlvn/versioned-go-control-plane/pkg/server/v3"
	"github.com/spf13/cobra"

	grpcserver1 "github.com/sefaphlvn/bigbang/grpc/grpcserver"
	grpcserver "github.com/sefaphlvn/bigbang/grpc/server"
	"github.com/sefaphlvn/bigbang/grpc/server/bridge"
	"github.com/sefaphlvn/bigbang/grpc/server/snapshot"
	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/log"
	"github.com/sefaphlvn/bigbang/pkg/models"
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
		appContext := db.NewMongoDB(appConfig, logger)
		ctxCache := snapshot.GetContext(logger)
		grpcserver1.ResetGrpcServerNodeIDs(appContext.Client)
		// go grpcserver1.ScheduleSetNodeIDs(ctxCache, db.Client)

		//pokeServer := poke.NewPokeServer(ctxCache, db, logger, appConfig)
		activeClients := &models.ActiveClients{Clients: make(map[string]*models.Client)}
		activeClientsService := bridge.NewActiveClientsService(logger)
		pokeService := bridge.NewPokeService(ctxCache, appContext)
		activeClientsService.ActiveClients = activeClients

		callbacks := grpcserver.NewCallbacks(pokeService, ctxCache, activeClientsService, appContext)
		srv := server.NewServer(context.Background(), ctxCache.Cache.Cache, callbacks)
		grpcServer := grpcserver.NewServer(srv, port, logger, ctxCache, activeClients)

		grpcServer.Run(appContext)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.PersistentFlags().UintVar(&port, "port", 18000, "xDS management server port")
	grpcCmd.PersistentFlags().StringVar(&nodeID, "nodeID", "test", "Node ID")
}
