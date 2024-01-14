package config_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aronkst/go-rate-limit/internal/config"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("ENV_TEST", "true")

	t.Setenv("REDIS_ADDRESS", "localhost")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "0")
	t.Setenv("DEFAULT_IP_MAX_REQ_PER_SEC", "5")
	t.Setenv("DEFAULT_TOKEN_MAX_REQ_PER_SEC", "10")
	t.Setenv("DEFAULT_IP_BLOCK_DURATION", "15s")
	t.Setenv("DEFAULT_TOKEN_BLOCK_DURATION", "30s")
	t.Setenv("CUSTOM_MAX_REQ_PER_SEC", "127.0.0.1=100;abc123=200")
	t.Setenv("CUSTOM_BLOCK_DURATION", "127.0.0.1=10s;abc123=20s")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.RedisAddress != "localhost" {
		t.Errorf("RedisAddress = %v, want %v", cfg.RedisAddress, "localhost")
	}

	if cfg.RedisPassword != "secret" {
		t.Errorf("RedisPassword = %v, want %v", cfg.RedisPassword, "secret")
	}

	if cfg.RedisDB != 0 {
		t.Errorf("RedisDB = %v, want %v", cfg.RedisDB, 0)
	}

	if cfg.DefaultIPMaxReqPerSec != 5 {
		t.Errorf("DefaultIPMaxReqPerSec = %v, want %v", cfg.DefaultIPMaxReqPerSec, 5)
	}

	if cfg.DefaultTokenMaxReqPerSec != 10 {
		t.Errorf("DefaultTokenMaxReqPerSec = %v, want %v", cfg.DefaultTokenMaxReqPerSec, 10)
	}

	if cfg.DefaultIPBlockDuration != 15*time.Second {
		t.Errorf("DefaultIPBlockDuration = %v, want %v", cfg.DefaultIPBlockDuration, 15*time.Second)
	}

	if cfg.DefaultTokenBlockDuration != 30*time.Second {
		t.Errorf("DefaultTokenBlockDuration = %v, want %v", cfg.DefaultTokenBlockDuration, 30*time.Second)
	}

	expectedCustomMaxReqPerSec := map[string]int{"127.0.0.1": 100, "abc123": 200}
	if !reflect.DeepEqual(cfg.CustomMaxReqPerSec, expectedCustomMaxReqPerSec) {
		t.Errorf("CustomMaxReqPerSec = %v, want %v", cfg.CustomMaxReqPerSec, expectedCustomMaxReqPerSec)
	}

	expectedCustomBlockDuration := map[string]time.Duration{"127.0.0.1": 10 * time.Second, "abc123": 20 * time.Second}
	if !reflect.DeepEqual(cfg.CustomBlockDuration, expectedCustomBlockDuration) {
		t.Errorf("CustomBlockDuration = %v, want %v", cfg.CustomBlockDuration, expectedCustomBlockDuration)
	}
}
