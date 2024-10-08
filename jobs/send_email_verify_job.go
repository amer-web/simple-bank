package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const SendEmail = "send_email"

type SendEmailVerifyJob struct {
	Username string
	Email    string
	Token    string
}

func (rD *RedisDistributor) DistributorSendEmailToQueue(ctx context.Context,
	payload SendEmailVerifyJob, opts ...asynq.Option) error {
	jsonPayload, _ := json.Marshal(payload)
	task := asynq.NewTask(SendEmail, jsonPayload, opts...)
	info, err := rD.client.EnqueueContext(ctx, task)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}
	log.Info().Str("type", info.Type).Msg("enqueued task ")
	return nil
}

func (rD *RedisProcessor) ProcessEmailVerify(ctx context.Context, task *asynq.Task) error {
	var payload SendEmailVerifyJob
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	subject := "verify Email"
	content := `<p>Verify Email</p>`
	content += fmt.Sprintf("<h4>welcome : %s</h4>", payload.Username)
	verifyLink := fmt.Sprintf("<a href=\"%s/v1/verify_email?token=%s\">Verify Email</a>", "http://localhost:8080", payload.Token)
	content += fmt.Sprintf("<p>To verify your email, please click %s</p>", verifyLink)

	to := []string{payload.Email}
	err := rD.sender.SendEmail(subject, content, to, nil)
	if err != nil {
		return err
	}
	log.Info().Str("type", task.Type()).Msg("task processed")
	return nil
}
