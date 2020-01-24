package cfg

import "testing"

func validToml() string {
	return `
favicon = "static/img/favicon.png"
refresh = 10
title = "My Health Monitor"

[[checks]]
service = "Lua Service"
description = "Always set healthy"
lua = "setHealthy()"

[[modules]]
name = "greetings"
lua = """
local M = {}

function M.sayHello()
  print('Why hello there')
end

return M
"""
`
}

func TestLoadConfig_Valid(t *testing.T) {
	config, err := LoadConfig(validToml())
	if err != nil {
		t.Fatal("Should not return an error")
	}

	if config == nil {
		t.Errorf("Config should not be nil")
	}
}

func TestLoadConfig_Invalid(t *testing.T) {
	invalidToml := "-=$!$#bzrg43"
	config, err := LoadConfig(invalidToml)
	if err == nil {
		t.Errorf("Config should return an error")
	}

	if config != nil {
		t.Errorf("Config should be nil")
	}
}
