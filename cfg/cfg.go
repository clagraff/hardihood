package cfg

import (
	"github.com/BurntSushi/toml"
	"github.com/clagraff/hardihood/checkup"
	"github.com/clagraff/hardihood/service"
)

type Config interface {
	Title() string
	Favicon() string
	Refresh() int
	Services() []service.Service
	Scripts() []struct {
		Name string
		Lua  string
	}
}

type config struct {
	title    string
	favicon  string
	refresh  int
	services []service.Service
	scripts  []struct {
		Name string
		Lua  string
	}
}

func (c config) Title() string               { return c.title }
func (c config) Favicon() string             { return c.favicon }
func (c config) Refresh() int                { return c.refresh }
func (c config) Services() []service.Service { return c.services }
func (c config) Scripts() []struct {
	Name string
	Lua  string
} {
	return c.scripts
}

type tomlConfig struct {
	Title   string
	Favicon string
	Refresh int
	Checks  []struct {
		Service     string
		Description string
		Lua         string
	}
	Scripts []struct {
		Name string
		Lua  string
	}
}

func (c tomlConfig) Services() []service.Service {
	serviceNames := []string{}
	serviceSet := make(map[string][]checkup.Checkup)

	for _, cfgCheck := range c.Checks {
		serviceName := cfgCheck.Service
		check := checkup.Make(cfgCheck.Description, cfgCheck.Lua)

		if _, ok := serviceSet[serviceName]; !ok {
			serviceSet[serviceName] = []checkup.Checkup{check}
			serviceNames = append(serviceNames, serviceName)
		} else {
			serviceSet[serviceName] = append(serviceSet[serviceName], check)
		}
	}

	allServices := []service.Service{}
	for _, serviceName := range serviceNames {
		checks := serviceSet[serviceName]
		allServices = append(
			allServices,
			service.Make(serviceName, checks),
		)
	}

	return allServices
}

func (c tomlConfig) ToConfig() config {
	return config{
		title:    c.Title,
		favicon:  c.Favicon,
		refresh:  c.Refresh,
		services: c.Services(),
		scripts:  c.Scripts,
	}
}

func LoadConfig(data string) (Config, error) {
	var cfg Config

	tomlCfg := tomlConfig{}
	_, err := toml.Decode(data, &tomlCfg)
	if err != nil {
		return cfg, err
	}
	return tomlCfg.ToConfig(), nil
}
