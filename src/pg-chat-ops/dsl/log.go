package dsl

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func dslLog(level string, L *lua.LState) {
	log.Printf("[%s] %s\n", level, L.CheckAny(1).String())
}

func (d *dslConfig) dslLogError(L *lua.LState) int {
	dslLog("ERROR", L)
	return 0
}

func (d *dslConfig) dslLogInfo(L *lua.LState) int {
	dslLog("INFO", L)
	return 0
}
