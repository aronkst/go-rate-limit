package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	RedisAddress              string
	RedisPassword             string
	RedisDB                   int
	DefaultIPMaxReqPerSec     int
	DefaultTokenMaxReqPerSec  int
	DefaultIPBlockDuration    time.Duration
	DefaultTokenBlockDuration time.Duration
	CustomMaxReqPerSec        map[string]int
	CustomBlockDuration       map[string]time.Duration
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		RedisAddress:              viper.GetString("REDIS_ADDRESS"),
		RedisPassword:             viper.GetString("REDIS_PASSWORD"),
		RedisDB:                   viper.GetInt("REDIS_DB"),
		DefaultIPMaxReqPerSec:     viper.GetInt("DEFAULT_IP_MAX_REQ_PER_SEC"),
		DefaultTokenMaxReqPerSec:  viper.GetInt("DEFAULT_TOKEN_MAX_REQ_PER_SEC"),
		DefaultIPBlockDuration:    viper.GetDuration("DEFAULT_IP_BLOCK_DURATION"),
		DefaultTokenBlockDuration: viper.GetDuration("DEFAULT_TOKEN_BLOCK_DURATION"),
	}

	cfg.CustomMaxReqPerSec, err = parseCustomMaxReqPerSec(viper.GetString("CUSTOM_MAX_REQ_PER_SEC"))
	if err != nil {
		return nil, err
	}

	cfg.CustomBlockDuration, err = parseCustomBlockDuration(viper.GetString("CUSTOM_BLOCK_DURATION"))
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseCustomMaxReqPerSec(customSettings string) (map[string]int, error) {
	result := make(map[string]int)

	pairs := strings.Split(customSettings, ";")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			maxReqPerSec, err := stringToInt(kv[1])
			if err != nil {
				return nil, err
			}

			result[kv[0]] = maxReqPerSec
		}
	}

	return result, nil
}

func parseCustomBlockDuration(customSettings string) (map[string]time.Duration, error) {
	result := make(map[string]time.Duration)

	pairs := strings.Split(customSettings, ";")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			blockDuration, err := stringToDuration(kv[1])
			if err != nil {
				return nil, err
			}

			result[kv[0]] = blockDuration
		}
	}

	return result, nil
}

func stringToInt(value string) (int, error) {
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func stringToDuration(value string) (time.Duration, error) {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}

	return duration, nil
}
