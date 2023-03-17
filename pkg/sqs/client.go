package queue

import (
	"context"
	"poc-ddb-tidb-search/pkg/session"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var (
	sqsClient *sqs.Client
)

func NewSQSClient(ctx context.Context) *sqs.Client {
	if sqsClient != nil {
		return sqsClient
	}

	cfg, success := session.GetSessionConfig(ctx)
	if !success {
		return nil
	}

	sqsClient = sqs.NewFromConfig(cfg)

	return sqsClient
}
