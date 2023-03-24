package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/sirupsen/logrus"

	"poc-ddb-tidb-search/pkg/db"
	"poc-ddb-tidb-search/pkg/models"

	"github.com/google/uuid"
)

const (
	recType_Insert = "insert"
	recType_Modify = "modify"
	recType_Remove = "remove"
)

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
}

func handler(ctx context.Context, event events.SQSEvent) (*events.SQSEventResponse, error) {

	tiDB, err := db.NewTiDB(db.TiDB_DatabaseName)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to connect to TiDB instance")
	}
	defer tiDB.Close()

	failures := make([]events.SQSBatchItemFailure, 0, len(event.Records))
	for _, record := range event.Records {
		err := parseAndSend(ctx, &record, tiDB)
		if err != nil {
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
		}
	}

	return &events.SQSEventResponse{
		BatchItemFailures: failures,
	}, nil
}

func parseAndSend(ctx context.Context, record *events.SQSMessage, tiDB db.DB) error {

	job := new(models.Job)
	err := json.Unmarshal([]byte(record.Body), job)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "SQSParseErr",
			"body":  record.Body,
		}).Error("failed to parse sqs message")
		return err
	}

	// get sqs messag attributes -> recordType and orgID
	orgID, recordType := getMessageAttribs(record)

	return syncToTiDB(ctx, orgID, recordType, job, tiDB)
}

func syncToTiDB(ctx context.Context, orgID, recordType string, job *models.Job, tiDB db.DB) error {
	var err error

	switch recordType {

	case recType_Insert:
		err = insertToTiDB(job, tiDB)

	case recType_Modify:
		//_, err = tiDB.Put(job) // deffered

	case recType_Remove:
		//err = tiDB.Delete(nil) // deffered

	default:
		return errors.New("unknown record type")
	}

	return err
}

func getMessageAttribs(record *events.SQSMessage) (string, string) {
	orgID, recType := "", ""

	for k, v := range record.MessageAttributes {
		switch k {
		case "orgID":
			orgID = *v.StringValue
		case "recordType":
			recType = *v.StringValue
		default:
			continue
		}
	}

	return orgID, recType
}

func insertToTiDB(job *models.Job, tiDB db.DB) error {
	sqlStmts := make([]string, 0)
	id := uuid.New()

	jobStmt, err := db.MakeInsertJobSQLStatement(job, id.String())
	if err != nil {
		return err
	}

	jobRefStmt, err := db.MakeInsertJobReferenceSQLStatement(job, id.String())
	if err != nil {
		return err
	}

	sqlStmts = append(sqlStmts, jobStmt, jobRefStmt)

	_, err = tiDB.Put(sqlStmts[0], sqlStmts[1])
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
