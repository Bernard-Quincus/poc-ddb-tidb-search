package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"poc-ddb-tidb-search/pkg/models"
	queue "poc-ddb-tidb-search/pkg/sqs"

	"poc-ddb-tidb-search/pkg/ddbstream"

	logger "github.com/sirupsen/logrus"
)

var (
	sqsClient queue.SendMessageClient
	QueueURL  = ""
)

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
	QueueURL = os.Getenv("QUEUE_URL")
}

func handler(ctx context.Context, streamEvent events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	if err := initSQSClient(ctx); err != nil {
		return events.DynamoDBEventResponse{
			BatchItemFailures: []events.DynamoDBBatchItemFailure{{ItemIdentifier: ""}},
		}, err
	}

	streamEventFailures := events.DynamoDBEventResponse{BatchItemFailures: make([]events.DynamoDBBatchItemFailure, 0)}

	for _, event := range streamEvent.Records {
		job, err := parseRecordStream(&event)
		if err != nil {
			logger.WithFields(logger.Fields{
				"error": err.Error(),
				"code":  "DDBStreamErr",
			}).Error("failed to parse stream event")
			streamEventFailures.BatchItemFailures = append(streamEventFailures.BatchItemFailures, events.DynamoDBBatchItemFailure{ItemIdentifier: event.Change.SequenceNumber})
			continue
		}

		err = sendToQueue(ctx, job, event.EventName)
		if err != nil {
			logger.WithFields(logger.Fields{
				"error": err.Error(),
				"code":  "SQSErr",
			}).Error("failed to send stream event to queue ")
			streamEventFailures.BatchItemFailures = append(streamEventFailures.BatchItemFailures, events.DynamoDBBatchItemFailure{ItemIdentifier: event.Change.SequenceNumber})
		}
	}

	return streamEventFailures, nil
}

func parseRecordStream(event *events.DynamoDBEventRecord) (*models.Job, error) {
	job := new(models.Job)
	err := ddbstream.ConvertStreamRecord(&event.Change, "json", job) // should add dynamodbav tags in job struct as well
	if err != nil {
		return nil, err
	}

	return job, nil
}

func sendToQueue(ctx context.Context, job *models.Job, recordType string) error {
	msg, err := json.Marshal(job)
	if err != nil {
		return err
	}
	input := queue.NewQueueInput(QueueURL, string(msg))
	input.SetMessageAttributes("recordType", strings.ToLower(recordType)) // remove, insert or modify
	input.SetMessageAttributes("orgID", job.OrgID2)
	input.SetMessageGroupID(job.OrgID2)

	_, err = queue.Enqueue(ctx, sqsClient, input)
	if err != nil {
		return err
	}

	return nil
}

func initSQSClient(ctx context.Context) error {
	if sqsClient != nil {
		return nil
	}

	client := queue.NewSQSClient(ctx)
	if client == nil {
		logger.WithFields(logger.Fields{
			"error": "failed to initialise sqs client",
			"code":  "SQSErr",
		}).Error("failed to parse request body")
		return errors.New("failed to initialise sqs client")
	}

	sqsClient = client
	return nil
}

func main() {
	lambda.Start(handler)
}
