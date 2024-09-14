package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/sefaphlvn/bigbang/grpc/poke"
	grpcserver "github.com/sefaphlvn/bigbang/grpc/server"
	"github.com/sefaphlvn/bigbang/pkg/config"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/log"
	"github.com/spf13/cobra"

	_ "net/http/pprof"
)

var (
	port   uint
	nodeID string
)

func stats() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	go func() {
		for {
			helper.PrintMemoryUsage()
			helper.PrintCPUUsage()
			time.Sleep(3 * time.Second)
		}
	}()
}

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "server-grpc",
	Short: "Start Bigbang GRPC Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		var appConfig = config.Read(cfgFile)
		var logger = log.NewLogger(appConfig)
		var db = db.NewMongoDB(appConfig, logger)

		var ctxCache = grpcserver.GetContext(logger)

		var pokeServer = poke.NewPokeServer(ctxCache, db, logger, appConfig)
		go pokeServer.Run(pokeServer)

		var callbacks = grpcserver.NewCallbacks(logger)
		var srv = server.NewServer(context.Background(), ctxCache.Cache.Cache, callbacks)
		var grpcServer = grpcserver.NewServer(srv, port, logger, ctxCache)

		grpcServer.Run()
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.PersistentFlags().UintVar(&port, "port", 18000, "xDS management server port")
	grpcCmd.PersistentFlags().StringVar(&nodeID, "nodeID", "test", "Node ID")
}
