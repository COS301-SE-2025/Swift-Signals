package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Port  int    `env:"PORT"  envDefault:"8080"`
	Env   string `env:"ENV"   envDefault:"development"`
	Debug bool   `env:"DEBUG" envDefault:"false"`
	NoTag string // Should be ignored
}

func TestLoadWithEnvVariables(t *testing.T) {
	setEnv(t, "PORT", "9000")
	setEnv(t, "ENV", "production")
	setEnv(t, "DEBUG", "true")

	var cfg TestConfig
	err := Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, 9000, cfg.Port)
	assert.Equal(t, "production", cfg.Env)
	assert.Equal(t, true, cfg.Debug)
}

func TestLoadWithDefaults(t *testing.T) {
	unsetEnv(t, "PORT")
	unsetEnv(t, "ENV")
	unsetEnv(t, "DEBUG")

	var cfg TestConfig
	err := Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, false, cfg.Debug)
}

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set env %s: %v", key, err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv(key); err != nil {
			t.Logf("Failed to unset env %s: %v", key, err)
		}
	})
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Failed to unset env %s: %v", key, err)
	}
}
