package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"poc-ddb-tidb-search/pkg/models"
	"poc-ddb-tidb-search/pkg/query"

	"poc-ddb-tidb-search/pkg/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/sirupsen/logrus"
)

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger.Info("reading params")

	var start = time.Now()

	params, err := query.ParametersFromRequest(&request)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "ParamsErr",
		}).Error("failed to parse request parameters")
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	orgID := query.GetOrgID(&request)

	tiDB, err := db.NewTiDB(db.TiDB_DatabaseName)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to connect to TiDB instance")
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	defer tiDB.Close()

	res, err := searchInTiDB(tiDB, params, orgID, start)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to query TiDB")
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	res.PageSize = params.PageSize

	jsonRes, err := json.Marshal(res)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "JSONErr",
		}).Error("failed to marshal jobs")
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonRes),
	}, nil
}

func searchInTiDB(tiDB db.DB, params *query.JobSearchParams, orgID string, start time.Time) (*query.JobSearchResult, error) {

	sqlStms := db.MakeSearchSQLStatements(params, orgID)

	// temp
	logger.Info("sqlStms:", sqlStms)

	res, err := tiDB.Search(sqlStms[0], sqlStms[1])
	if err != nil {
		return nil, err
	}

	return makeFinalResult(res, start)
}

func makeFinalResult(res any, start time.Time) (*query.JobSearchResult, error) {
	var dbResult = res.(*db.TiDBResult)
	var result = new(query.JobSearchResult)
	result.Data = make([]*query.JobRow, 0)
	result.TotalItems = dbResult.TotalItems

	for _, row := range dbResult.Details {
		jobJson, _ := base64.StdEncoding.DecodeString(row.Detail)
		var job = new(models.Job)
		json.Unmarshal(jobJson, job)
		result.Data = append(result.Data, &query.JobRow{UUID: row.UUID, Job: job})
	}

	result.ResponseTime = time.Since(start).String()
	return result, nil
}

func main() {
	lambda.Start(handler)
}
