package main

import (
	"context"
	"fmt"
	"net"

	pb "go-app/server/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedServiceServerServer
}

func (s *server) SayHi(ctx context.Context, req *pb.CliRequest) (*pb.SerResponse, error) {
	return &pb.SerResponse{
		ResponseMsg: "hoo",
	}, nil
}

func main() {

	r := gin.Default()

	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterServiceServerServer(grpcServer, &server{})

	go func() {
		err := grpcServer.Serve(listen)
		if err != nil {
			fmt.Printf("failed to serve: %v", err)
		}
	}()

	r.GET("/api2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"massage": "hi, greetting from server 2",
		})
	})

	r.Run(":8082")
}
