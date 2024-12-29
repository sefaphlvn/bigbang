package bridge

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"

	"github.com/sefaphlvn/bigbang/pkg/bridge"
	"github.com/sefaphlvn/bigbang/pkg/db"
)

type AppHandler struct {
	Context          *db.AppContext
	GRPCConn         *grpc.ClientConn
	SnapshotResource bridge.SnapshotResourceServiceClient
	SnapshotKeys     bridge.SnapshotKeyServiceClient
	Poke             bridge.PokeServiceClient
	ActiveClients    bridge.ActiveClientsServiceClient
}

func ipv4Dialer(ctx context.Context, addr string) (net.Conn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp4", addr)
	if err != nil {
		conn, err = d.DialContext(ctx, "tcp6", addr)
	}
	return conn, err
}

func NewBridgeHandler(appCtx *db.AppContext) *AppHandler {
	var transportCredentials credentials.TransportCredentials
	if appCtx.Config.BigbangTLSEnabled == "true" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		transportCredentials = credentials.NewTLS(tlsConfig)
	} else {
		transportCredentials = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(
		appCtx.Config.BigbangAddress+":"+appCtx.Config.BigbangPort,
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithContextDialer(ipv4Dialer),
		grpc.WithDisableServiceConfig(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithAuthority(appCtx.Config.BigbangAddress),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	/* if err := checkHealth(conn); err != nil {
		log.Fatalf("gRPC server health check failed: %v", err)
	} */

	return &AppHandler{
		Context:          appCtx,
		GRPCConn:         conn,
		SnapshotResource: bridge.NewSnapshotResourceServiceClient(conn),
		SnapshotKeys:     bridge.NewSnapshotKeyServiceClient(conn),
		Poke:             bridge.NewPokeServiceClient(conn),
		ActiveClients:    bridge.NewActiveClientsServiceClient(conn),
	}
}

func checkHealth(conn *grpc.ClientConn) error {
	healthClient := grpc_health_v1.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startTime := time.Now()
	for {
		resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: ""})
		if err == nil && resp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
			fmt.Printf("Health check failed: %v. Retrying...\n", err)
			return nil
		}

		fmt.Printf("Health check failed: %v. Retrying...\n", err)
		if time.Since(startTime) >= 30*time.Second {
			break
		}

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("Health check failed after retries within 30 seconds")
}

func (h *AppHandler) Close() {
	if h.GRPCConn != nil {
		h.GRPCConn.Close()
	}
}
