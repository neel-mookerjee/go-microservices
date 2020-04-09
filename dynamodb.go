package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DbInterface interface {
	GetPendingDownsamplePreviewItems() ([]DownsamplingObject, error)
	UpdateDownsamplePreviewItem(DownsampleObject DownsamplingObject) (string, error)
	GetDownsamplingItem(queryId string) (DownsamplingObject, error)
	GetMetricsStacks() (MetricsStacks, error)
}

func (d *Dynamodb) formTableName(table string) *string {
	return aws.String(fmt.Sprintf("%s%s", d.Configs.DbTablePrefix, table))
}

type Dynamodb struct {
	Configs *Config
	Svc     dynamodbiface.DynamoDBAPI
}

func NewDynamodb(config *Config) (*Dynamodb, error) {
	dbsession, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		return nil, err
	}

	return &Dynamodb{Svc: dynamodb.New(dbsession), Configs: config}, nil
}

func (d *Dynamodb) GetDownsamplingItem(queryId string) (DownsamplingObject, error) {
	input := &dynamodb.GetItemInput{
		TableName: d.formTableName("metrics_downsample_queries"),
		Key: map[string]*dynamodb.AttributeValue{
			"queryId": {
				S: aws.String(queryId),
			},
		},
	}
	var query DownsamplingObject
	result, err := d.Svc.GetItem(input)
	err = dynamodbattribute.UnmarshalMap(result.Item, &query)
	if err != nil {
		return query, err
	}
	return query, nil
}

func (d *Dynamodb) GetPendingDownsamplePreviewItems() ([]DownsamplingObject, error) {
	return d.getDownsamplePreviewItems("PREVIEW_PENDING")
}

func (d *Dynamodb) getDownsamplePreviewItems(state string) ([]DownsamplingObject, error) {
	input := &dynamodb.ScanInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":queryState": {
				S: aws.String(state),
			},
		},
		FilterExpression: aws.String("queryState = :queryState"),
		TableName:        d.formTableName("metrics_downsample_queries"),
	}
	var ds []DownsamplingObject
	result, err := d.Svc.Scan(input)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &ds)
	if err != nil {
		return ds, err
	}
	return ds, nil
}

func (d *Dynamodb) UpdateDownsamplePreviewItem(DownsampleObject DownsamplingObject) (string, error) {
	ds, err := dynamodbattribute.MarshalMap(DownsampleObject)
	if err != nil {
		return "Failure", err
	}

	_, err = d.Svc.PutItem(
		&dynamodb.PutItemInput{
			TableName: d.formTableName("metrics_downsample_queries"),
			Item:      ds,
		})
	if err != nil {
		return "Failure", err
	}

	status := "Success"
	return status, nil
}

func (d *Dynamodb) GetMetricsStacks() (MetricsStacks, error) {
	input := &dynamodb.ScanInput{
		TableName: d.formTableName("customer_metrics_stacks"),
	}
	var msList MetricsStacks
	result, err := d.Svc.Scan(input)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &msList)
	return msList, err
}
