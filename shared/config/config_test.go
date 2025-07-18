package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Port  int    `env:"PORT" envDefault:"8080"`
	Env   string `env:"ENV" envDefault:"development"`
	Debug bool   `env:"DEBUG" envDefault:"false"`
	NoTag string // Should be ignored
}

func TestLoadWithEnvVariables(t *testing.T) {
	os.Setenv("PORT", "9000")
	os.Setenv("ENV", "production")
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("ENV")
	defer os.Unsetenv("DEBUG")

	var cfg TestConfig
	err := Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, 9000, cfg.Port)
	assert.Equal(t, "production", cfg.Env)
	assert.Equal(t, true, cfg.Debug)
}

func TestLoadWithDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("ENV")
	os.Unsetenv("DEBUG")

	var cfg TestConfig
	err := Load(&cfg)
	assert.NoError(t, err)

	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, false, cfg.Debug)
}
