package cfg

import (
	"github.com/BurntSushi/toml"
	"github.com/clagraff/hardihood/checkup"
	"github.com/clagraff/hardihood/service"
)

// Config is used to provide the program with configurable properties.
type Config interface {
	// Title returns the title of the status site.
	Title() string

	// Favicon is used to provide a link or data url for the site's favicon.
	Favicon() string

	// Refresh is the number of seconds to wait before
	// auto-refreshing the page. Negative values will prevent any refreshing.
	Refresh() int

	// Services returns a list of all services to check.
	Services() []service.Service

	// Scripts is used to specific Lua modules which will be made
	// available to all service check scripts.
	Scripts() []struct {
		// Name of the lua module.
		Name string

		// Lua code.
		Lua string
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

// Title returns the title of the site.
func (c config) Title() string { return c.title }

// Favicon returns the path or data url for the site favicon.
func (c config) Favicon() string { return c.favicon }

// Refresh returns the number of seconds to wait before auto-refreshing.
func (c config) Refresh() int { return c.refresh }

// Services returns a list of services to check.
func (c config) Services() []service.Service { return c.services }

// Scripts returns a list of lua modules to provide to all status checks.
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

// Services parses all status checks to get the set of defined services.
func (c tomlConfig) services() []service.Service {
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

func (c tomlConfig) toConfig() config {
	return config{
		title:    c.Title,
		favicon:  c.Favicon,
		refresh:  c.Refresh,
		services: c.services(),
		scripts:  c.Scripts,
	}
}

// LoadConfig takes a toml config and returns a Config.
func LoadConfig(data string) (Config, error) {
	var cfg Config

	tomlCfg := tomlConfig{}
	_, err := toml.Decode(data, &tomlCfg)
	if err != nil {
		return cfg, err
	}
	return tomlCfg.toConfig(), nil
}
