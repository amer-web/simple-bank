package middleware

import (
	"context"
	"errors"
	tok "github.com/amer-web/simple-bank/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

var protectedMethods = map[string]bool{
	"/pb.SimpleBank/LoginUser":  false,
	"/pb.SimpleBank/CreateUser": true,
	"/pb.SimpleBank/UpdateUser": true,
}
var protectedMethodsApi = map[string]bool{
	"/v1/create_user": true,
	"/v1/update_user": true,
	"/v1/login":       false,
}

func AuthInterceptorGrpc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if secured, ok := protectedMethods[info.FullMethod]; secured && ok {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		token := md["authorization"]
		if len(token) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token not provided")
		}
		tokenStr := strings.TrimPrefix(token[0], "Bearer ")
		handleToken := tok.NewMakerToken()
		_, err := handleToken.VerifyToken(tokenStr)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}
	}
	return handler(ctx, req)
}
func AuthMiddlewareGrpcGateway(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if protectedMethodsApi[r.URL.Path] {
			token := r.Header.Get("Authorization")
			tokenStr := strings.TrimPrefix(token, "Bearer ")
			handleToken := tok.NewMakerToken()
			_, err := handleToken.VerifyToken(tokenStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
