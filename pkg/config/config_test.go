package config

import "testing"

func TestConfigInit(t *testing.T) {
	if Config.Region == "" {
		t.Error("Config.Region is empty")
		return
	}
	t.Logf("Config: %v", Config)
}
