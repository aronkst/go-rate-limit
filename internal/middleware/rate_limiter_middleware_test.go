package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aronkst/go-rate-limit/internal/middleware"
)

type MockRateLimit struct {
	Exceeded bool
	Error    error
}

func (m *MockRateLimit) IsLimitExceeded(identifier string) (bool, error) {
	return m.Exceeded, m.Error
}

func TestRateLimiterMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		exceeded         bool
		error            error
		expectedHTTPCode int
		expectedBody     string
	}{
		{"LimitNotExceeded", false, nil, http.StatusOK, ""},
		{"LimitExceeded", true, nil, http.StatusTooManyRequests, "you have reached the maximum number of requests or actions allowed within a certain time frame\n"},
		{"InternalServerError", false, errors.New("error"), http.StatusInternalServerError, "internal server error\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockLimiter := &MockRateLimit{
				Exceeded: tc.exceeded,
				Error:    tc.error,
			}

			middleware := middleware.RateLimiterMiddleware(mockLimiter)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			testHandler := middleware(handler)

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			testHandler.ServeHTTP(w, req)

			if w.Code != tc.expectedHTTPCode {
				t.Errorf("Test %s failed: Expected HTTP status code %d, got %d", tc.name, tc.expectedHTTPCode, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("Test %s failed: Expected HTTP body %s, got %s", tc.name, tc.expectedBody, w.Body.String())
			}
		})
	}
}
