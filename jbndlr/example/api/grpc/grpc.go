package grpc

import (
	"context"
	"fmt"
	"jbndlr/example/api"
	"jbndlr/example/api/grpc/pb"
	"jbndlr/example/conf"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	// Status : This API's status information
	Status = api.NewStatus()
)

type gRPCServer struct {
	pb.UnimplementedExampleServer
}

// Serve : Start serving gRPC API.
func Serve(port int16) error {
	server := grpc.NewServer()
	pb.RegisterExampleServer(server, &gRPCServer{})
	reflection.Register(server)

	start := func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return err
		}

		Status.Start()
		log.Printf("Serving gRPC :%d\n", port)

		return server.Serve(lis)
	}

	stop := func() error {
		log.Printf("Stopping gRPC")
		stopped := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(stopped)
		}()

		timer := time.NewTimer(time.Duration(conf.P.API.GracefulSeconds) * time.Second)
		select {
		case <-timer.C:
			server.Stop()
			return fmt.Errorf("Forced shutdown after timeout")
		case <-stopped:
			timer.Stop()
		}

		return nil
	}

	err := api.ServeGracefully(start, stop)
	Status.Stop(err)
	return err
}

// Info : implements function for gRPCServer.
func (s *gRPCServer) Info(ctx context.Context, in *pb.Empty) (*pb.InfoMessage, error) {
	return &pb.InfoMessage{
		Service:  conf.P.Service.Name,
		Version:  conf.P.Service.Version,
		GrpcPort: int32(conf.P.API.GRPCPort),
		HttpPort: int32(conf.P.API.RESTPort),
	}, nil
}
