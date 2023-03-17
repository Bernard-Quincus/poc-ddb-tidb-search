package db

import (
	"context"
	"poc-ddb-tidb-search/pkg/session"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	client *dynamodb.Client
)

func getDDBClient(ctx context.Context) *dynamodb.Client {
	if client != nil {
		return client
	}

	cfg, success := session.GetSessionConfig(ctx)
	if !success {
		return nil
	}

	client = dynamodb.NewFromConfig(cfg)
	return client
}
