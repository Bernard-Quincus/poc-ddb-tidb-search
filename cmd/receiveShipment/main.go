package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"poc-ddb-tidb-search/pkg/db"
	"poc-ddb-tidb-search/pkg/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/sirupsen/logrus"
)

var (
	ddb       db.DB
	tableName = ""
)

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
	tableName = os.Getenv("POC_TABLE")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	job, err := parseBody(req.Body)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "ParseErr",
		}).Error("failed to parse request body")
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	ddb = db.NewDDB(ctx)
	err = ddb.SetTableName(tableName)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TablenameErr",
		}).Error("invalid table name")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	err = insertJob(job)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "CFGErr",
		}).Error("failed to insert record into DynamoDB")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func parseBody(body string) (*models.Job, error) {
	job := new(models.Job)

	err := json.Unmarshal([]byte(body), job)
	if err != nil {
		return nil, err
	}

	job.OrgID2 = "POC-TEST-ORGID-001122" // pk
	job.DocID = job.ID                   // sk
	return job, nil
}

func insertJob(job *models.Job) error {
	_, err := ddb.Put(job)
	return err
}

func main() {
	lambda.Start(handler)
}
