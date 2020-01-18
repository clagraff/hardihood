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
	"github.com/kohkimakimoto/gluayaml"
	"github.com/yuin/gluare"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

func getCSS() string {
	return `
		body {
			background-color: #fdfdfd;
			text-align: center;
			font-family: Sans;
		}

		section {
			display: block;
			background: #ffffff;
			border: 1px solid #F0F0F0;
			border-radius: 9px;
			text-align: left;
			margin: 2rem auto 2rem auto;
			width: 960px;
			padding: 1rem;

			-webkit-box-shadow: 0px 0px 5px 0px rgba(0,0,0,0.30);
			-moz-box-shadow: 0px 0px 5px 0px rgba(0,0,0,0.30);
			box-shadow: 0px 0px 5px 0px rgba(0,0,0,0.30);
		}

		section header {
			font-size: 2rem;
		}

		section ul {
			list-style: none;
			padding: 0px;
			margin: 0 1rem;
		}

		section li {
			border: 1px solid #F0F0F0;
			border-radius: 9px;	
			background-color: #FBFBFB;
			padding: 1rem;
			margin: 1rem 0 1rem 0;
		}

		section li span {
			float: right;
		}

		span.healthy {
			color: green;
		}

		span.sick {
			color: red;
		}

		span.healthy .icon {
			background-color: green;
		}

		span.sick .icon {
			background-color: red;
		}

		span.icon {
			border-radius: 4rem;
			width: 1.2rem;
			text-align: center;
			margin-left: 1rem;
			color: white;
			font-weight: bold;
		}

`
}

func listingHTML() string {
	return `
<html>
	<head>
		<title>Status Page</title>
		<meta http-equiv="refresh" content="15">
		<link rel="icon" 
			type="image/png" 
			href="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiA/PjxzdmcgZW5hYmxlLWJhY2tncm91bmQ9Im5ldyAwIDAgNjQgNjQiIGhlaWdodD0iNjRweCIgdmVyc2lvbj0iMS4xIiB2aWV3Qm94PSIwIDAgNjQgNjQiIHdpZHRoPSI2NHB4IiB4bWw6c3BhY2U9InByZXNlcnZlIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIj48ZyBpZD0iTGF5ZXJfMSI+PGc+PGNpcmNsZSBjeD0iMzIiIGN5PSIzMiIgZmlsbD0iI0M3NUM1QyIgcj0iMzIiLz48L2c+PGcgb3BhY2l0eT0iMC4yIj48Zz48cGF0aCBkPSJNNDkuOTgyLDMxLjAwM2MtMC4wOTQtNS41MjItNC41NzQtMTAuNDQyLTEwLjEwNy0xMC40NDJjLTMuMiwwLTYuMDE5LDEuNjc0LTcuODc1LDQuMTMxICAgICBjLTEuODU2LTIuNDU3LTQuNjc2LTQuMTMxLTcuODc1LTQuMTMxYy01LjUzMywwLTEwLjAxMiw0LjkyMS0xMC4xMDcsMTAuNDQySDE0YzAsMC4wMzQsMC4wMDcsMC4wNjUsMC4wMDcsMC4wOTkgICAgIGMwLDAuMDI1LTAuMDA3LDAuMDQ5LTAuMDA3LDAuMDc2YzAsMC4xNTUsMC4wMzgsMC4yNzIsMC4wNDUsMC40MjFjMC40OTUsMTQuMDcxLDE3LjgxMywxOS44NCwxNy44MTMsMTkuODQgICAgIHMxNy41NzItNS43NjIsMTguMDkyLTE5LjgxOEM0OS45NTksMzEuNDY0LDUwLDMxLjM0LDUwLDMxLjE3OGMwLTAuMDI3LTAuMDA3LTAuMDUyLTAuMDA3LTAuMDc2YzAtMC4wMzYsMC4wMDctMC4wNjUsMC4wMDctMC4wOTkgICAgIEg0OS45ODJ6IiBmaWxsPSIjMjMxRjIwIi8+PC9nPjwvZz48Zz48Zz48cGF0aCBkPSJNNDkuOTgyLDI5LjAwM2MtMC4wOTQtNS41MjItNC41NzQtMTAuNDQyLTEwLjEwNy0xMC40NDJjLTMuMiwwLTYuMDE5LDEuNjc0LTcuODc1LDQuMTMxICAgICBjLTEuODU2LTIuNDU3LTQuNjc2LTQuMTMxLTcuODc1LTQuMTMxYy01LjUzMywwLTEwLjAxMiw0LjkyMS0xMC4xMDcsMTAuNDQySDE0YzAsMC4wMzQsMC4wMDcsMC4wNjUsMC4wMDcsMC4wOTkgICAgIGMwLDAuMDI1LTAuMDA3LDAuMDQ5LTAuMDA3LDAuMDc2YzAsMC4xNTUsMC4wMzgsMC4yNzIsMC4wNDUsMC40MjFjMC40OTUsMTQuMDcxLDE3LjgxMywxOS44NCwxNy44MTMsMTkuODQgICAgIHMxNy41NzItNS43NjIsMTguMDkyLTE5LjgxOEM0OS45NTksMjkuNDY0LDUwLDI5LjM0LDUwLDI5LjE3OGMwLTAuMDI3LTAuMDA3LTAuMDUyLTAuMDA3LTAuMDc2YzAtMC4wMzYsMC4wMDctMC4wNjUsMC4wMDctMC4wOTkgICAgIEg0OS45ODJ6IiBmaWxsPSIjRkZGRkZGIi8+PC9nPjwvZz48L2c+PGcgaWQ9IkxheWVyXzIiLz48L3N2Zz4=">
		<style>
			{{.CSS}}
		</style>
	</head>
	<body>
	{{range .Services}}
		<section>
			<header>{{.Name}}</header>
			<ul>
				{{range .Checks}}
				<li>
					{{.Description}}
					{{$status := .Status}}
					<span class="{{$status.CSSIdent}}">{{$status.Name}}<span class="icon">{{$status.HTMLChar}}</span></span>
				</li>
				{{end}}
			</ul>
		</section>
	{{end}}
	</body>
</html>
`
}

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
	err := state.DoString(c.luaScript)
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
	Checks []struct {
		Service     string
		Description string
		Lua         string
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
		}{
			CSS:      getCSS(),
			Services: cfg.Services(),
		}

		html := listingHTML()

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
