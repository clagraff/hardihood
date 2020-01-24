package scripting

import (
	"fmt"
	"net/http"

	"github.com/ailncode/gluaxmlpath"
	"github.com/cjoudrey/gluahttp"
	"github.com/cjoudrey/gluaurl"
	"github.com/clagraff/hardihood/status"
	"github.com/kohkimakimoto/gluayaml"
	"github.com/yuin/gluare"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type statusResult struct {
	isHealthy bool
}

func (l *statusResult) SetHealthy(_ *lua.LState) int {
	l.isHealthy = true
	return 1
}

func (l *statusResult) SetSick(_ *lua.LState) int {
	l.isHealthy = false
	return 1
}

type Script interface {
	Name() string
	Code() string
}

type script struct {
	name string
	code string
}

func (s script) Name() string { return s.name }
func (s script) Code() string { return s.code }

func MakeScript(name, code string) Script {
	return script{
		name: name,
		code: code,
	}
}

type preload struct {
	name string
	fn   *lua.LGFunction
}

type luaScript struct {
	name string
	fn   *lua.LFunction
}

type State interface {
	Execute(string) status.Status
	Preload(string, lua.LGFunction)
	Script(string, *lua.LFunction)
}

type state struct {
	preloads []preload
	scripts  []luaScript
}

func (s *state) Preload(name string, fn *lua.LGFunction) {
	s.preloads = append(
		s.preloads,
		preload{
			name: name,
			fn:   fn,
		},
	)
}

func (s *state) Script(name string, fn *lua.LFunction) {
	s.scripts = append(
		s.scripts,
		luaScript{
			name: name,
			fn:   fn,
		},
	)
}

func Execute(code string) status.Status {
	result := new(statusResult)

	state := lua.NewState()
	defer state.Close()
	state.SetGlobal("setHealthy", state.NewFunction(result.SetHealthy))
	state.SetGlobal("setSick", state.NewFunction(result.SetSick))
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
		return status.Sick
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

	err = state.DoString(code)
	if err != nil {
		fmt.Println(code)
		fmt.Println(err)
		return status.Sick
	}

	if result.isHealthy {
		return status.Healthy
	}

	return status.Sick
}
