package store

import (
	"fmt"

	"github.com/google/uuid"
)

// Datastore is a generic interface for a datastore
type Datastore interface {
	List(castTo interface{}, page, limit int64, options ...interface{}) error
	Store(item interface{}) error
	Get(key uuid.UUID, castTo interface{}) error
	Update(key uuid.UUID, castTo interface{}, changeSet ChangeSet) error
	Delete(key uuid.UUID) error
}

// NotFoundError error type for records not found in the datastore
type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Item not found")
}

// ChangeSet is a generic interface to map attribute changes
type ChangeSet map[string]interface{}
