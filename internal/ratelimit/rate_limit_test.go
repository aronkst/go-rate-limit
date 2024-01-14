package ratelimit_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aronkst/go-rate-limit/internal/config"
	"github.com/aronkst/go-rate-limit/internal/ratelimit"
)

type MockStorage struct {
	RequestCounts map[string]int
	Error         error
}

func (m *MockStorage) IncrementRequestCount(identifier string, expiry time.Duration) (int, error) {
	if m.Error != nil {
		return 0, m.Error
	}

	m.RequestCounts[identifier]++
	return m.RequestCounts[identifier], nil
}

func (m *MockStorage) GetRequestCount(identifier string) (int, error) {
	if m.Error != nil {
		return 0, m.Error
	}

	count, exists := m.RequestCounts[identifier]

	if !exists {
		return 0, nil
	}

	return count, nil
}

func TestIsLimitExceeded(t *testing.T) {
	cfg := &config.Config{
		DefaultIPMaxReqPerSec:     5,
		DefaultTokenMaxReqPerSec:  10,
		DefaultIPBlockDuration:    time.Minute,
		DefaultTokenBlockDuration: 2 * time.Minute,
		CustomMaxReqPerSec:        make(map[string]int),
		CustomBlockDuration:       make(map[string]time.Duration),
	}
	cfg.CustomMaxReqPerSec["testToken"] = 1

	mockStore := &MockStorage{RequestCounts: make(map[string]int)}

	rl := ratelimit.NewRateLimit(cfg, mockStore)

	tests := []struct {
		name         string
		identifier   string
		increment    int
		expectExceed bool
		expectError  bool
	}{
		{"IP NotExceeded", "192.168.1.1", 4, false, false},
		{"IP Exceeded", "192.168.1.1", 5, true, false},
		{"Token NotExceeded", "testToken", 0, false, false},
		{"Token Exceeded", "testToken", 1, true, false},
		{"StorageError", "errorToken", 0, false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore.RequestCounts[tc.identifier] = tc.increment
			if tc.expectError {
				mockStore.Error = errors.New("storage error")
			} else {
				mockStore.Error = nil
			}

			exceeded, err := rl.IsLimitExceeded(tc.identifier)

			if tc.expectError {
				if err == nil {
					t.Errorf("%s: expected an error but got none", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("%s: unexpected error: %v", tc.name, err)
				}
			}

			if exceeded != tc.expectExceed {
				t.Errorf("%s: expected exceeded to be %v, got %v", tc.name, tc.expectExceed, exceeded)
			}
		})
	}
}
