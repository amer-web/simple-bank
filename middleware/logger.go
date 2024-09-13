package middleware

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func LoggerInterceptorGrpc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	timeNow := time.Now()
	re, err := handler(ctx, req)
	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	duration := time.Since(timeNow)
	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}
	logger.Dur("duration", duration).
		Str("method", info.FullMethod).
		Int("statusCode", int(statusCode)).
		Str("status_text", statusCode.String()).
		Str("protocol", "grpc").
		Msgf("%s", "received a grpc request")

	return re, err
}

type ResponseWriterNew struct {
	statusCode int
	http.ResponseWriter
	body []byte
}

func (r *ResponseWriterNew) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
func (r *ResponseWriterNew) Write(body []byte) (int, error) {
	r.body = body
	return r.ResponseWriter.Write(body)
}

func LoggerInterceptorHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeNow := time.Now()
		rec := &ResponseWriterNew{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(rec, r)
		duration := time.Since(timeNow)
		logger := log.Info()
		if rec.statusCode >= 400 && rec.statusCode < 500 {
			logger = log.Error().Bytes("body", rec.body)
		}

		logger.
			Dur("duration", duration).
			Str("method", r.Method).
			Int("statusCode", rec.statusCode).
			Str("status_text", http.StatusText(rec.statusCode)).
			Str("path", r.RequestURI).
			Str("protocol", "http").
			Msgf("%s", "received a http request")

	})
}
