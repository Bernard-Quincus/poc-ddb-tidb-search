package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// implement me

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "test",
	}, nil
}

func main() {
	lambda.Start(handler)
}
