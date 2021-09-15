package store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

// CreateConnection to dynamodb
func CreateConnection(endpoint string) (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(endpoint)},
	)
	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}

// DynamoDBTable is a reference to a specific table in DynamoDB
type DynamoDBTable struct {
	table string
	conn  *dynamodb.DynamoDB
}

// NewDynamoDBTable creates a new reference to a DynamoDB table
func NewDynamoDBTable(conn *dynamodb.DynamoDB, table string) *DynamoDBTable {
	return &DynamoDBTable{
		conn: conn, table: table,
	}
}

// List gets a collection of resources
func (dt *DynamoDBTable) List(castTo interface{}) error {
	results, err := dt.conn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(dt.table),
	})
	if err != nil {
		return err
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &castTo); err != nil {
		return err
	}
	return nil
}

// Store a new Item
func (dt *DynamoDBTable) Store(item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dt.table),
	}
	_, err = dt.conn.PutItem(input)
	if err != nil {
		return err
	}
	return err
}

// Get an item
func (dt *DynamoDBTable) Get(key uuid.UUID, castTo interface{}) error {

	id, err := key.MarshalBinary()
	if err != nil {
		return err
	}

	result, err := dt.conn.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dt.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				B: id,
			},
		},
	})
	if err != nil {
		return err
	}
	if result.Item == nil {
		return &NotFoundError{}
	}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &castTo); err != nil {
		return err
	}
	return nil
}

// Update an item
func (dt *DynamoDBTable) Update(key uuid.UUID, attributes map[string]*dynamodb.AttributeValue, expression string) error {
	var err error

	id, err := key.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = dt.conn.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(dt.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				B: id,
			},
		},
		ExpressionAttributeValues: attributes,
		UpdateExpression:          aws.String(expression),
	})
	if err != nil {
		return err
	}
	// if err := dynamodbattribute.UnmarshalMap(result.Item, &castTo); err != nil {
	// 	return err
	// }
	return nil
}
