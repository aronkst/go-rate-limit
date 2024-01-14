package storage

import (
	"time"
)

type StorageInterface interface {
	IncrementRequestCount(identifier string, expiry time.Duration) (int, error)
	GetRequestCount(identifier string) (int, error)
}
