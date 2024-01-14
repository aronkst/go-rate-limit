package ratelimit

import (
	"regexp"
	"time"

	"github.com/aronkst/go-rate-limit/internal/config"
	"github.com/aronkst/go-rate-limit/internal/storage"
)

type RateLimit struct {
	config *config.Config
	Store  storage.StorageInterface
}

func NewRateLimit(config *config.Config, store storage.StorageInterface) *RateLimit {
	return &RateLimit{
		config: config,
		Store:  store,
	}
}

func (rl *RateLimit) IsLimitExceeded(identifier string) (bool, error) {
	currentCount, err := rl.Store.GetRequestCount(identifier)
	if err != nil {
		return false, err
	}

	maxReqPerSec := getMaxReqPerSec(identifier, rl.config)
	blockDuration := getBlockDuration(identifier, rl.config)

	if currentCount >= maxReqPerSec {
		return true, nil
	}

	_, err = rl.Store.IncrementRequestCount(identifier, blockDuration)

	if err != nil {
		return false, err
	}

	return false, nil
}

func isToken(identifier string) bool {
	tokenRegex := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
	return tokenRegex.MatchString(identifier)
}

func getMaxReqPerSec(identifier string, config *config.Config) int {
	if maxReqPerSec, ok := config.CustomMaxReqPerSec[identifier]; ok {
		return maxReqPerSec
	}

	if isToken(identifier) {
		return config.DefaultTokenMaxReqPerSec
	}

	return config.DefaultIPMaxReqPerSec
}

func getBlockDuration(identifier string, config *config.Config) time.Duration {
	if blockDuration, ok := config.CustomBlockDuration[identifier]; ok {
		return blockDuration
	}

	if isToken(identifier) {
		return config.DefaultTokenBlockDuration
	}

	return config.DefaultIPBlockDuration
}
