package database

import (
	"acortadorUrlService/url-api/model"
	"context"
	"errors"
	"fmt"
	"time"

	"acortadorUrlService/components/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DDBClient struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDDBClient(ctx context.Context, cfg *config.AppConfig) (*DDBClient, error) {
	awsCfg := aws.NewConfig().WithRegion(cfg.Region)
	sess, err := session.NewSession(awsCfg)
	if err != nil {
		return nil, err
	}

	table := cfg.TableName
	if table == "" {
		return nil, errors.New("DDB_TABLE (TableName) not set in config")
	}

	return &DDBClient{
		client:    dynamodb.New(sess),
		tableName: table,
	}, nil
}

func (d *DDBClient) SaveURL(ctx context.Context, hash string, original string) error {
	entry := model.ShortenedURL{
		Hash:      hash,
		Original:  original,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	item, err := dynamodbattribute.MarshalMap(entry)
	if err != nil {
		return err
	}

	_, err = d.client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	})
	return err
}

func (d *DDBClient) GetURL(ctx context.Context, hash string) (string, error) {
	res, err := d.client.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"hash": {S: aws.String(hash)},
		},
	})
	if err != nil {
		return "", err
	}
	if res.Item == nil {
		return "", fmt.Errorf("the shorten url was not found")
	}

	var entry model.ShortenedURL
	if err := dynamodbattribute.UnmarshalMap(res.Item, &entry); err != nil {
		return "", err
	}
	return entry.Original, nil
}

func (d *DDBClient) DeleteURL(ctx context.Context, hash string) error {
	_, err := d.client.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"hash": {S: aws.String(hash)},
		},
	})
	return err
}
