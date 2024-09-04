package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/amer-web/simple-bank/api"
	"github.com/amer-web/simple-bank/config"
	db2 "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/gapi"
	"github.com/amer-web/simple-bank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	source := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.Source.DBUser, config.Source.DBPassword, config.Source.DBHost, config.Source.DBPort, config.Source.DBName)
	db, err := sql.Open(config.Source.DRIVER, source)
	if err != nil {
		log.Fatal("error opening db:", err.Error())
	}
	defer db.Close()
	store := db2.NewStore(db)
	go runGrpcGateway(store)
	runGrpcServer(store)
}
func runGinServer(store db2.Store) {
	server := api.NewServer(store)
	server.Run()
}
func runGrpcServer(store db2.Store) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := gapi.NewServer(store)
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterSimpleBankServer(grpcServer, server)
	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func runGrpcGateway(store db2.Store) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	server := gapi.NewServer(store)

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	err := pb.RegisterSimpleBankHandlerServer(ctx, mux, server)
	if err != nil {
		return err
	}
	log.Println("HTTP server listening on :8080")
	return http.ListenAndServe(":8080", mux)
}
