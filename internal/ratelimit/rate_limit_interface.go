package ratelimit

type RateLimitInterface interface {
	IsLimitExceeded(identifier string) (bool, error)
}
