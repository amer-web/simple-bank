package jobs

import (
	"context"
	"github.com/hibiken/asynq"
)

type Distributor interface {
	DistributorSendEmailToQueue(ctx context.Context,
		payload SendEmailVerifyJob, opts ...asynq.Option) error
}

type RedisDistributor struct {
	client *asynq.Client
}

func NewRedisDistributor(opt asynq.RedisClientOpt) Distributor {
	client := asynq.NewClient(opt)
	return &RedisDistributor{
		client: client,
	}
}
