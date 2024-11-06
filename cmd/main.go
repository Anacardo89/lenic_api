package main

import (
	"log"
	"net"

	"github.com/Anacardo89/lenic_api/config"
	"github.com/Anacardo89/lenic_api/internal/endpoints"
	"github.com/Anacardo89/lenic_api/internal/pb"
	"github.com/Anacardo89/lenic_api/internal/server"
	"github.com/Anacardo89/lenic_api/pkg/db"
	"github.com/Anacardo89/lenic_api/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	logger.CreateLogger()
	logger.Info.Println("System start")

	// DB
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		logger.Error.Fatalln("Could not load dbConfig:", err)
	}
	db.Dbase, err = db.LoginDB(dbConfig)
	if err != nil {
		logger.Error.Fatalln("Could not connect to DB: ", err)
	}
	logger.Info.Println("Connecting to DB OK")

	// Server
	server.Server, err = config.LoadServerConfig()
	if err != nil {
		logger.Error.Fatalln("Could not load serverConfig:", err)
	}
	logger.Info.Println("Loading serverConfig OK")

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	pb.RegisterAuthServiceServer(s, &endpoints.AuthService{})
	pb.RegisterLenicServer(s, &endpoints.ApiService{})

	lis, err := net.Listen("tcp", ":"+server.Server.GrpcPort)
	if err != nil {
		logger.Error.Fatalln("failed to listen: ", err)
	}

	log.Printf("Starting gRPC listener on port " + server.Server.GrpcPort)
	if err := s.Serve(lis); err != nil {
		logger.Error.Fatalln("failed to serve: ", err)
	}
}
