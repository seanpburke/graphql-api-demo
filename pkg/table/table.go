package table

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/sburke-at-ziprecruiter/graphql-api-demo/pkg/config"
)

// Item is an interface for items that go in the table.
type Item interface {
	MakePK() string // Constructs the item's partition (HASH) key.
	MakeSK() string // Constructs the item's sort (RANGE) key.
	Init()          // sets item.PK and item.SK
}

func GetItem(item Item) error {
	result, err := config.Config.DDB.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(config.Config.Table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(item.MakePK()),
			},
			"SK": {
				S: aws.String(item.MakeSK()),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("table.GetItem %w", err)
	}
	return dynamodbattribute.UnmarshalMap(result.Item, item)
}

func PutItem(i Item) error {
	i.Init()
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return fmt.Errorf("customer.PutItem dynamodbattribute.MarshalMap failed, %w", err)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(config.Config.Table),
	}
	if _, err = config.Config.DDB.PutItem(input); err != nil {
		return fmt.Errorf("dynamodb.PutItem(%s) failed, %w ", config.Config.Table, err)
	}
	return nil
}

func CreateTable() (*dynamodb.CreateTableOutput, error) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("GSI2PK"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{ // GSI1 is a reverse index PK -> SK
				IndexName: aws.String("GSI1"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("SK"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("PK"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String(dynamodb.ProjectionTypeAll),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
			},
			{ // GSI2 is a secondary index GSI2PK -> PK
				IndexName: aws.String("GSI2"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("GSI2PK"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("PK"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String(dynamodb.ProjectionTypeAll),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(config.Config.Table),
	}

	output, err := config.Config.DDB.CreateTable(input)
	if err != nil {
		err = fmt.Errorf("CreateTable(%q) %w", config.Config.Table, err)
	}
	return output, err
}
