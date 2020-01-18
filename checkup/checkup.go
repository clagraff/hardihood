package checkup

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

type Checkup interface {
	Description() string
	Status() status.Status
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

type luaCheckup struct {
	description string
	luaScript   string
}

func (c luaCheckup) Description() string { return c.description }
func (c luaCheckup) Status() status.Status {
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

	err = state.DoString(c.luaScript)
	if err != nil {
		fmt.Println(c.luaScript)
		fmt.Println(err)
		return status.Sick
	}

	if result.isHealthy {
		return status.Healthy
	}

	return status.Sick
}

func MakeCheckup(desc string, luaScript string) Checkup {
	return luaCheckup{
		description: desc,
		luaScript:   luaScript,
	}
}
