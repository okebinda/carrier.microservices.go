package store

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
func (dt *DynamoDBTable) List(castTo interface{}, page, limit int64) error {

	// get first page of results
	currentPage := int64(1)
	results, err := dt.conn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(dt.table),
		Limit:     aws.Int64(limit),
	})
	if err != nil {
		return err
	}

	// support for pagination: beyond page 1
	for currentPage < page && len(results.LastEvaluatedKey) > 0 {
		results, err = dt.conn.Scan(&dynamodb.ScanInput{
			TableName:         aws.String(dt.table),
			Limit:             aws.Int64(limit),
			ExclusiveStartKey: results.LastEvaluatedKey,
		})
		if err != nil {
			return err
		}
		currentPage++
	}

	// populate output with results
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

	fmt.Println(fmt.Sprintf("Item: %v", result.Item))

	return nil
}

// Update an item
func (dt *DynamoDBTable) Update(key uuid.UUID, castTo interface{}, changeSet ChangeSet) error {
	var err error

	// get binary value of ID
	id, err := key.MarshalBinary()
	if err != nil {
		return err
	}

	// initialize vars that dictate update behavior
	updateAttributes := map[string]*dynamodb.AttributeValue{}
	updateExpressions := []string{}

	// initialize vars that dictate remove behavior
	removeAttributes := []string{}

	// loop over change set and build update vars
	for k, v := range changeSet {
		placeholder := fmt.Sprintf(":%s", k)
		switch v.(type) {
		case string:
			updateAttributes[placeholder] = &dynamodb.AttributeValue{
				S: aws.String(v.(string)),
			}
			updateExpressions = append(updateExpressions, fmt.Sprintf("%s=:%s", k, k))
		case int:
			updateAttributes[placeholder] = &dynamodb.AttributeValue{
				N: aws.String(strconv.Itoa(v.(int))),
			}
			updateExpressions = append(updateExpressions, fmt.Sprintf("%s=:%s", k, k))
		case int64:
			updateAttributes[placeholder] = &dynamodb.AttributeValue{
				N: aws.String(strconv.FormatInt(v.(int64), 10)),
			}
			updateExpressions = append(updateExpressions, fmt.Sprintf("%s=:%s", k, k))
		case bool:
			updateAttributes[placeholder] = &dynamodb.AttributeValue{
				BOOL: aws.Bool(v.(bool)),
			}
			updateExpressions = append(updateExpressions, fmt.Sprintf("%s=:%s", k, k))
		case []string:
			updateAttributes[placeholder] = &dynamodb.AttributeValue{
				SS: aws.StringSlice(v.([]string)),
			}
			updateExpressions = append(updateExpressions, fmt.Sprintf("%s=:%s", k, k))
		case map[string]string:
			val, err := dynamodbattribute.MarshalMap(v.(map[string]string))
			if err != nil {
				return err
			}
			updateAttributes[placeholder] = &dynamodb.AttributeValue{
				M: val,
			}
			updateExpressions = append(updateExpressions, fmt.Sprintf("%s=:%s", k, k))
		case time.Time:
			val := v.(time.Time)
			if val.IsZero() {
				// remove nil times
				// why? sparse indexing: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-indexes-general-sparse-indexes.html
				// also nil times don't make sense
				removeAttributes = append(removeAttributes, k)
			} else {
				updateAttributes[placeholder] = &dynamodb.AttributeValue{
					S: aws.String(val.Format("2006-01-02T15:04:05Z07:00")),
				}
				updateExpressions = append(updateExpressions, fmt.Sprintf("%s=:%s", k, k))
			}
		}
	}

	udpateExpression := fmt.Sprintf("SET %s", strings.Join(updateExpressions, ", "))
	if len(removeAttributes) > 0 {
		udpateExpression = udpateExpression + fmt.Sprintf(" REMOVE %s", strings.Join(removeAttributes, ", "))
	}

	// perform update
	result, err := dt.conn.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(dt.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				B: id,
			},
		},
		ExpressionAttributeValues: updateAttributes,
		UpdateExpression:          aws.String(udpateExpression),
		ReturnValues:              aws.String("ALL_NEW"),
	})
	if err != nil {
		return err
	}

	// update original object with updated values
	if err = dynamodbattribute.UnmarshalMap(result.Attributes, &castTo); err != nil {
		return err
	}

	return nil
}

// Delete an item
func (dt *DynamoDBTable) Delete(key uuid.UUID) error {

	id, err := key.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = dt.conn.DeleteItem(&dynamodb.DeleteItemInput{
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

	return nil
}
