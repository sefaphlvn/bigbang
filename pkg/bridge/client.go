package bridge

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/sefaphlvn/bigbang/pkg/db"
)

func ipv4Dialer(ctx context.Context, addr string) (net.Conn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp4", addr)
	if err != nil {
		conn, err = d.DialContext(ctx, "tcp6", addr)
	}
	return conn, err
}

func NewGRPCClient(appCtx *db.AppContext) (*grpc.ClientConn, error) {
	var transportCredentials credentials.TransportCredentials
	if appCtx.Config.BigbangTLSEnabled == "true" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		transportCredentials = credentials.NewTLS(tlsConfig)
	} else {
		transportCredentials = insecure.NewCredentials()
	}

	return grpc.NewClient(
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
} 