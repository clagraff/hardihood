package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/clagraff/hardihood/checkup"
	"github.com/clagraff/hardihood/service"
	"github.com/clagraff/hardihood/static"
)

type config struct {
	Title   string
	Favicon string
	Refresh int
	Checks  []struct {
		Service     string
		Description string
		Lua         string
	}
	Scripts []struct {
		Lua string
	}
}

func (c config) Services() []service.Service {
	serviceNames := []string{}
	serviceSet := make(map[string][]checkup.Checkup)

	for _, cfgCheck := range c.Checks {
		serviceName := cfgCheck.Service
		check := checkup.MakeCheckup(cfgCheck.Description, cfgCheck.Lua)

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
			service.MakeService(serviceName, checks),
		)
	}

	return allServices
}

func LoadConfig(data string) (config, error) {
	cfg := config{}
	_, err := toml.Decode(data, &cfg)
	return cfg, err
}

func main() {
	contents, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	cfg, err := LoadConfig(string(contents))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := struct {
			CSS      string
			Services []service.Service
			Config   config
		}{
			CSS:      static.DefaultCSS(),
			Services: cfg.Services(),
			Config:   cfg,
		}

		html := static.HTML()

		tmpl, err := template.New("page").Parse(html)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		err = tmpl.Execute(w, page)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
