package dsl

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

const (
	tgTocken = `697096392:AAE6v7MbYE9_p82rAmPC2HQEykdEGnY6MUU`
	tgChatID = `80734283`
)

func testTelgram(t *testing.T) {
	tg := &dslTelegram{
		tocken:    tgTocken,
		ignoreSsl: true,
		proxy:     "http://192.168.184.28:3128",
	}
	if err := tg.buildHttpClient(); err != nil {
		t.Fatalf("build client: %s\n", err.Error())
	}
	message := `
привет пидор усатый
    `
	if err := tg.sendTextMessage(tgChatID, message); err != nil {
		t.Fatalf("send message: %s\n", err.Error())
	}

	fd, err := os.Open(filepath.Join("tests", "test.png"))
	if err != nil {
		t.Fatalf("send message: %s\n", err.Error())
	}

	st, _ := fd.Stat()
	log.Printf("size: %d\n", st.Size())

	message = `
крутая титла
    `

	if err := tg.sendPhotoWithMessage(tgChatID, message, "markdown", fd); err != nil {
		t.Fatalf("send message: %s\n", err.Error())
	}
}

func TestTelgramPhoto(t *testing.T) {

	state := lua.NewState()
	Register(NewConfig(), state)
	if err := state.DoFile(filepath.Join("tests", "telegram_test.lua")); err != nil {
		t.Fatalf("execute lua error: %s\n", err.Error())
	}
}
