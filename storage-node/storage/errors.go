package storage

import "fmt"

type DuplicateKeyError struct {
	Message string
	Key     string
}

func (e *DuplicateKeyError) Error() string {
	return e.Message
}

type KeyNotFoundError struct {
	Message string
	Key     string
}

func (e *KeyNotFoundError) Error() string {
	return e.Message
}

func NewKeyNotFoundError(key string) *KeyNotFoundError {
	return &KeyNotFoundError{
		Message: fmt.Sprintf("Key not found: %s", key),
		Key:     key,
	}
}
