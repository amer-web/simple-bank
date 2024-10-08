package jobs

import (
	"context"
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/mail"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskProcessor interface {
	Start() error
	ProcessEmailVerify(ctx context.Context, task *asynq.Task) error
}

type RedisProcessor struct {
	server *asynq.Server
	store  db.Store
	sender mail.EmailSender
}

func NewRedisProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Err(err).Str("type", task.Type()).Str("Payload", string(task.Payload())).
				Msg("can't process task")
		}),
		Logger: NewLogger(),
	})
	sender := mail.NewGmailSender()
	return &RedisProcessor{
		server: server,
		store:  store,
		sender: sender,
	}
}

func (p *RedisProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(SendEmail, p.ProcessEmailVerify)
	return p.server.Start(mux)
}
