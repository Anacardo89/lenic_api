package main

import (
	"crypto/tls"

	"github.com/Anacardo89/lenic_api/config"
	"github.com/Anacardo89/lenic_api/internal/pb"
	"github.com/Anacardo89/lenic_api/pkg/db"
	"github.com/Anacardo89/lenic_api/pkg/fsops"
	"github.com/Anacardo89/lenic_api/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	logger.CreateLogger()
	logger.Info.Println("System start")
	fsops.MakeImgDir()

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

	// Certificate
	cert, err := tls.LoadX509KeyPair(fsops.SSLCertificate, fsops.SSLkey)
	if err != nil {
		logger.Error.Fatalln("failed to load key pair: ", err)
	}
	logger.Info.Println("Loading SSL Certificates OK")

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	s := grpc.NewServer(opts...)

	pb.RegisterLenicServer(s, &server{})
}
