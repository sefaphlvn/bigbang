package suubar

import (
	"context"
	"log"
	"net"

	pb "github.com/sefaphlvn/bigbang/busproto"

	"google.golang.org/grpc"
)

// server is used to implement bootstrap.BootstrapServiceServer.
type server struct {
	pb.UnimplementedBootstrapServiceServer
}

// Bootstrap implements bootstrap.BootstrapServiceServer.
func (s *server) Bootstrap(_ context.Context, in *pb.BootstrapRequest) (*pb.BootstrapResponse, error) {
	log.Printf("Received: %v", in.GetDomain())
	return &pb.BootstrapResponse{Message: "Bootstrap successful for domain: " + in.GetDomain()}, nil
}

func Start() {
	lis, err := net.Listen("tcp", "localhost:50041")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBootstrapServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
