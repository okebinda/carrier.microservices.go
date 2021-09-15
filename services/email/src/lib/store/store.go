package store

import (
	"fmt"

	"github.com/google/uuid"
)

// Datastore is a generic interface for a datastore
type Datastore interface {
	List(castTo interface{}) error
	Get(key uuid.UUID, castTo interface{}) error
	Store(item interface{}) error
}

// NotFoundError error type for records not found in the datastore
type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Item not found")
}
