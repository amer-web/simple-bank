package main

import (
	"context"
	"fmt"
	"github.com/amer-web/simple-bank/api"
	"github.com/amer-web/simple-bank/config"
	db2 "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/gapi"
	"github.com/amer-web/simple-bank/jobs"
	"github.com/amer-web/simple-bank/middleware"
	"github.com/amer-web/simple-bank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
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
	err := config.LoadConfig(".")
	if err != nil {
		log.Error().Msg(err.Error())
	}
}
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	source := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.Source.DBUser, config.Source.DBPassword, config.Source.DBHost, config.Source.DBPort, config.Source.DBName)

	//connectionPool, err := pgxpool.New(context.Background(), source)

	config, err := pgxpool.ParseConfig(source)
	config.MinConns = 10

	connectionPoold, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error().Msgf("error opening db: %s", err.Error())

	}

	store := db2.NewStore(connectionPoold)
	redisConf := asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
		//Password: "yourpassword",
	}

	taskDistrebutor := jobs.NewRedisDistributor(redisConf)

	go runTasksProcessor(redisConf, store)

	go runGrpcGateway(store, taskDistrebutor)

	runGrpcServer(store, taskDistrebutor)
}
func runGinServer(store db2.Store) {
	server := api.NewServer(store)
	server.Run()
}
func runTasksProcessor(opt asynq.RedisClientOpt, store db2.Store) {
	server := jobs.NewRedisProcessor(opt, store)
	log.Info().Msg("start process tasks")
	server.Start()
}
func runGrpcServer(store db2.Store, taskDistrebutor jobs.Distributor) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Error().Msgf("failed to listen: %s", err.Error())

	}
	server := gapi.NewServer(store, taskDistrebutor)
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
func runGrpcGateway(store db2.Store, taskDistrebutor jobs.Distributor) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	server := gapi.NewServer(store, taskDistrebutor)

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080/"},
		AllowedMethods: []string{"POST"},
	})
	logger := middleware.AuthMiddlewareGrpcGateway(mux)
	log.Info().Msg("HTTP server listening on :8080")

	handler := c.Handler(middleware.LoggerInterceptorHttp(logger))

	httpServer := &http.Server{
		Handler: handler,
		Addr:    ":8080",
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil

}
