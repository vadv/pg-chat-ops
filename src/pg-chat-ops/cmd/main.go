package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	dsl "pg-chat-ops/dsl"

	lua "github.com/yuin/gopher-lua"
)

var (
	BuildVersion = "unknown"
	version      = flag.Bool("version", false, "print version and exit")
	initScript   = flag.String("init-script", "init.lua", "path to lua initial script")
)

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	if *version {
		fmt.Printf("%s\n", BuildVersion)
		os.Exit(1)
	}
	state := lua.NewState()
	dsl.Register(dsl.NewConfig(), state)
	if err := state.DoFile(*initScript); err != nil {
		log.Printf("[FATAL] execute %s: %s\n", *initScript, err.Error())
		os.Exit(2)
	}
}
