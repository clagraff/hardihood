package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/ailncode/gluaxmlpath"
	"github.com/cjoudrey/gluahttp"
	"github.com/cjoudrey/gluaurl"
	"github.com/clagraff/hardihood/static"
	"github.com/kohkimakimoto/gluayaml"
	"github.com/yuin/gluare"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type Status interface {
	Name() string
	CSSIdent() string
	HTMLChar() string
}

type status struct {
	name     string
	cssIdent string
	htmlChar string
}

func (s status) Name() string     { return s.name }
func (s status) CSSIdent() string { return s.cssIdent }
func (s status) HTMLChar() string { return s.htmlChar }

func MakeStatus(name, cssIdent, htmlChar string) Status {
	return status{
		name:     name,
		cssIdent: cssIdent,
		htmlChar: htmlChar,
	}
}

var Healthy Status = MakeStatus("Healthy", "healthy", "✓")
var Sick Status = MakeStatus("Sick", "sick", "✕")

type Check interface {
	Description() string
	Status() Status
}

type luaStatusResult struct {
	isHealthy bool
}

func (l *luaStatusResult) IsHealthy(_ *lua.LState) int {
	l.isHealthy = true
	return 1
}

func (l *luaStatusResult) IsSick(_ *lua.LState) int {
	l.isHealthy = false
	return 1
}

type luaCheck struct {
	description string
	luaScript   string
}

func (c luaCheck) Description() string { return c.description }
func (c luaCheck) Status() Status {
	result := new(luaStatusResult)

	state := lua.NewState()
	defer state.Close()
	state.SetGlobal("setHealthy", state.NewFunction(result.IsHealthy))
	state.SetGlobal("setSick", state.NewFunction(result.IsSick))
	state.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	state.PreloadModule("re", gluare.Loader)
	state.PreloadModule("json", luajson.Loader)
	state.PreloadModule("yaml", gluayaml.Loader)
	state.PreloadModule("url", gluaurl.Loader)
	state.PreloadModule("xmlpath", gluaxmlpath.Loader)

	greetings := `
local M = {}

local function sayMyName()
  print('general kenobi')
end

function M.sayHello()
  print('Why hello there')
  sayMyName()
end

return M
`
	fn, err := state.LoadString(greetings)
	if err != nil {
		fmt.Println(err)
		return Sick
	}

	state.SetField(
		state.GetField(
			state.GetField(
				state.Get(
					lua.EnvironIndex,
				),
				"package",
			),
			"preload",
		),
		"greetings",
		fn,
	)

	err = state.DoString(c.luaScript)
	if err != nil {
		fmt.Println(c.luaScript)
		fmt.Println(err)
		return Sick
	}

	if result.isHealthy {
		return Healthy
	}

	return Sick
}

func MakeCheck(desc string, luaScript string) luaCheck {
	return luaCheck{
		description: desc,
		luaScript:   luaScript,
	}
}

type Service interface {
	Name() string
	Checks() []Check
}

type service struct {
	name   string
	checks []luaCheck
}

func (s service) Name() string { return s.name }
func (s service) Checks() []Check {
	checks := []Check{}
	for _, c := range s.checks {
		checks = append(checks, c)
	}
	return checks
}

func MakeService(name string, checks []luaCheck) Service {
	return service{
		name:   name,
		checks: checks,
	}
}

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

func (c config) Services() []Service {
	serviceNames := []string{}
	services := make(map[string][]luaCheck)

	for _, cfgCheck := range c.Checks {
		serviceName := cfgCheck.Service
		check := MakeCheck(cfgCheck.Description, cfgCheck.Lua)

		if _, ok := services[serviceName]; !ok {
			services[serviceName] = []luaCheck{check}
			serviceNames = append(serviceNames, serviceName)
		} else {
			services[serviceName] = append(services[serviceName], check)
		}
	}

	allServices := []Service{}
	for _, serviceName := range serviceNames {
		checks := services[serviceName]
		allServices = append(
			allServices,
			MakeService(serviceName, checks),
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
			Services []Service
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
