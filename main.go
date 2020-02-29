package main

import (
	"context"
	"fmt"
	"github.com/dnk90/grpc_gorm_sample/sample"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {

	// there are 3 parts:
	// - connect to mysql through gorm
	// - start grpc server.
	// - start an grpc gateway to forward request to grpc server

	// load .env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	username := os.Getenv("MYSQL_USER_NAME")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName   := os.Getenv("MYSQL_DBNAME")
	host     := os.Getenv("MYSQL_HOST")
	serverEndpoint := os.Getenv("SERVER_ENDPOINT")
	gatewayEndpoint := os.Getenv("GATEWAY_ENDPOINT")

	// 1/ connect to mysql throug gorm
	db, err := gorm.Open("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8&parseTime=True&loc=Local", username, password, host, dbName))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	} ()

	if !db.HasTable(sample.ItemORM{}.TableName()) {
		// create new table
		db.CreateTable(sample.ItemORM{})
	}

	// 2/ Init and start server with above db
	lis, err := net.Listen("tcp", serverEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcSever := grpc.NewServer()
	sample.RegisterSampleServiceServer(grpcSever, sample.NewServer(db))
	go func() {
		if err = grpcSever.Serve(lis); err != nil {
			panic(err)
		}
	} ()

	// 3/ start grpc gateway

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err = sample.RegisterSampleServiceHandlerFromEndpoint(ctx, mux, serverEndpoint, opts)
	if err != nil {
		panic(err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	http.ListenAndServe(gatewayEndpoint, mux)
}
