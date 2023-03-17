package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	logger "github.com/sirupsen/logrus"
)

type dynamoDB struct {
	client    *dynamodb.Client
	ctx       context.Context
	tableName string
}

type ddbResult struct {
}

func NewDDB(ctx context.Context) DB {
	client := getDDBClient(ctx)
	return &dynamoDB{
		client: client,
		ctx:    ctx,
	}
}

func (ddb *dynamoDB) SetTableName(name string) error {
	if name == "" {
		return errors.New("empty table name")
	}
	ddb.tableName = name
	return nil
}

func (ddb *dynamoDB) Put(input any) (any, error) {
	if input == nil {
		return nil, nil
	}
	if ddb.tableName == "" {
		return nil, errors.New("no table specified")
	}

	av, err := attributevalue.MarshalMapWithOptions(input, func(opt *attributevalue.EncoderOptions) {
		opt.TagKey = "json" // should have dynamodbav tags so we dont need to pass options
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "DDBMarshalErr",
		}).Error("failed to marshal item for dynamodb")
		return nil, err
	}

	_, err = ddb.client.PutItem(ddb.ctx, &dynamodb.PutItemInput{
		TableName: aws.String(ddb.tableName),
		Item:      av,
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "DDBPutErr",
		}).Error("failed to put item in ddb")
		return nil, err
	}

	return nil, nil
}

func (ddb *dynamoDB) Get(input any) (any, error) {

	return nil, nil
}

func (ddb *dynamoDB) Delete(input any) error {

	return nil
}

func (ddb *dynamoDB) Search(input any) (any, error) {

	return nil, nil
}

func (ddb *dynamoDB) Close(input any) error {

	return nil
}
