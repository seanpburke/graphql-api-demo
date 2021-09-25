package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfigInit(t *testing.T) {

	assert.NotEqual(t, Config.AppName, "")
	assert.NotEqual(t, Config.Region, "")
	assert.NotEqual(t, Config.Registry, "")
	assert.NotEqual(t, Config.Table, "")
	assert.NotEqual(t, Config.Listen, "")

	// confirm env var overrides
	os.Setenv("LISTEN", "127.0.0.1:80")
	assert.Equal(t, viper.GetString("listen"), "127.0.0.1:80")

	t.Logf("Config: %v", Config)
}
