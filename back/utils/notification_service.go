package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/yanatoritakuma/budget/back/domain/budget"
)

type sqsNotificationService struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSNotificationService(ctx context.Context) (budget.INotificationService, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		return nil, fmt.Errorf("SQS_QUEUE_URL is not set")
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(cfg)
	return &sqsNotificationService{
		client:   client,
		queueURL: queueURL,
	}, nil
}

func (s *sqsNotificationService) SendBudgetExceededNotification(ctx context.Context, event budget.BudgetExceededEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &s.queueURL,
		MessageBody: &[]string{string(body)}[0],
	})
	return err
}
