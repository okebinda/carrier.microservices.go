package store

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

// Datastore is a generic interface for a datastore
type Datastore interface {
	List(castTo interface{}) error
	Get(key uuid.UUID, castTo interface{}) error
	Store(item interface{}) error
	Update(key uuid.UUID, attributes map[string]*dynamodb.AttributeValue, expression string) error
}

// NotFoundError error type for records not found in the datastore
type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Item not found")
}
