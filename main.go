package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/amer-web/simple-bank/api"
	"github.com/amer-web/simple-bank/config"
	db2 "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/gapi"
	"github.com/amer-web/simple-bank/middleware"
	"github.com/amer-web/simple-bank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		log.Error().Msg(err.Error())
	}
}
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	source := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.Source.DBUser, config.Source.DBPassword, config.Source.DBHost, config.Source.DBPort, config.Source.DBName)
	db, err := sql.Open(config.Source.DRIVER, source)
	if err != nil {
		log.Error().Msgf("error opening db: %s", err.Error())

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
		log.Error().Msgf("failed to listen: %s", err.Error())

	}
	server := gapi.NewServer(store)
	interceptors := grpc.ChainUnaryInterceptor(
		middleware.LoggerInterceptorGrpc,
		middleware.AuthInterceptorGrpc,
	)
	grpcServer := grpc.NewServer(interceptors)
	reflection.Register(grpcServer)
	pb.RegisterSimpleBankServer(grpcServer, server)
	log.Info().Msg("gRPC server listening on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Error().Msgf("failed to serve: %v", err.Error())
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

	logger := middleware.AuthMiddlewareGrpcGateway(mux)
	log.Info().Msg("HTTP server listening on :8080")
	err = http.ListenAndServe(":8080", middleware.LoggerInterceptorHttp(logger))
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil

}
