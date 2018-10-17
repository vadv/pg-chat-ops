package dsl

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// получение telegram из lua-state
func checkTelegram(L *lua.LState) *dslTelegram {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*dslTelegram); ok {
		return v
	}
	L.ArgError(1, "telegram expected")
	return nil
}

// создание telegram
func (c *dslConfig) dslNewTelegram(L *lua.LState) int {
	setValue := func(d *dslTelegram, t *lua.LTable, key string) {
		luaVal := t.RawGetString(key)
		if val, ok := luaVal.(lua.LString); ok {
			switch key {
			case "tocken":
				d.tocken = string(val)
			case "proxy":
				d.proxy = string(val)
			default:
				L.RaiseError("unknown option key: %s", key)
			}
		}
		if val, ok := luaVal.(lua.LBool); ok {
			switch key {
			case "ignore_ssl":
				d.ignoreSsl = bool(val)
			default:
				L.RaiseError("unknown option key: %s", key)
			}
		}
	}
	tocken := L.CheckString(1)
	d := newDSLTelegram(tocken)
	t := L.CheckTable(2)
	setValue(d, t, `proxy`)
	setValue(d, t, `ignore_ssl`)
	if err := d.buildHttpClient(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = d
	L.SetMetatable(ud, L.GetTypeMetatable("telegram"))
	L.Push(ud)
	log.Printf("[INFO] New %s\n", d.toString())
	return 1
}
