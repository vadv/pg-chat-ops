package dsl

import (
	"runtime"

	lua "github.com/yuin/gopher-lua"
)

type dslConfig struct{}

func NewConfig() *dslConfig {
	return &dslConfig{}
}

func Register(config *dslConfig, L *lua.LState) {

	plugin := L.NewTypeMetatable("plugin")
	L.SetGlobal("plugin", plugin)
	L.SetField(plugin, "new", L.NewFunction(config.dslNewPlugin))
	L.SetField(plugin, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"filename":    config.dslPluginFilename,
		"run":         config.dslPluginRun,
		"stop":        config.dslPluginStop,
		"error":       config.dslPluginError,
		"was_stopped": config.dslPluginWasStopped,
		"is_running":  config.dslPluginIsRunning,
	}))

	ioutil := L.NewTypeMetatable("ioutil")
	L.SetGlobal("ioutil", ioutil)
	L.SetField(ioutil, "readfile", L.NewFunction(config.dslIoutilReadFile))
	L.SetField(ioutil, "read_file", L.NewFunction(config.dslIoutilReadFile))

	filepath := L.NewTypeMetatable("filepath")
	L.SetGlobal("filepath", filepath)
	L.SetField(filepath, "base", L.NewFunction(config.dslFilepathBasename))
	L.SetField(filepath, "dir", L.NewFunction(config.dslFilepathDir))
	L.SetField(filepath, "ext", L.NewFunction(config.dslFilepathExt))
	L.SetField(filepath, "glob", L.NewFunction(config.dslFilepathGlob))

	os := L.NewTypeMetatable("goos")
	L.SetGlobal("goos", os)
	L.SetField(os, "stat", L.NewFunction(config.dslOsStat))
	L.SetField(os, "pagesize", L.NewFunction(config.dslOsPagesize))

	time := L.NewTypeMetatable("time")
	L.SetGlobal("time", time)
	L.SetField(time, "unix", L.NewFunction(config.dslTimeUnix))
	L.SetField(time, "unix_nano", L.NewFunction(config.dslTimeUnixNano))
	L.SetField(time, "sleep", L.NewFunction(config.dslTimeSleep))
	L.SetField(time, "parse", L.NewFunction(config.dslTimeParse))

	http := L.NewTypeMetatable("http")
	L.SetGlobal("http", http)
	L.SetField(http, "request", L.NewFunction(config.dslHttpRequest))
	L.SetField(http, "escape", L.NewFunction(config.dslHttpEscape))
	L.SetField(http, "unescape", L.NewFunction(config.dslHttpUnEscape))

	strings := L.NewTypeMetatable("strings")
	L.SetGlobal("strings", strings)
	L.SetField(strings, "split", L.NewFunction(config.dslStringsSplit))
	L.SetField(strings, "has_prefix", L.NewFunction(config.dslStringsHasPrefix))
	L.SetField(strings, "has_suffix", L.NewFunction(config.dslStringsHasSuffix))
	L.SetField(strings, "trim", L.NewFunction(config.dslStringsTrim))
	L.SetField(strings, "trim_prefix", L.NewFunction(config.dslStringsTrimPrefix))

	log := L.NewTypeMetatable("log")
	L.SetGlobal("log", log)
	L.SetField(log, "error", L.NewFunction(config.dslLogError))
	L.SetField(log, "info", L.NewFunction(config.dslLogInfo))

	crypto := L.NewTypeMetatable("crypto")
	L.SetGlobal("crypto", crypto)
	L.SetField(crypto, "md5", L.NewFunction(config.dslCryptoMD5))

	cmd := L.NewTypeMetatable("cmd")
	L.SetGlobal("cmd", cmd)
	L.SetField(cmd, "exec", L.NewFunction(config.dslCmdExec))

	goruntime := L.NewTypeMetatable("goruntime")
	L.SetGlobal("goruntime", goruntime)
	L.SetField(goruntime, "goarch", lua.LString(runtime.GOARCH))
	L.SetField(goruntime, "goos", lua.LString(runtime.GOOS))

	telegram := L.NewTypeMetatable("telegram")
	L.SetGlobal("telegram", telegram)
	L.SetField(telegram, "new", L.NewFunction(config.dslNewTelegram))
	L.SetField(telegram, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"send_text_message":     config.dslTelegramSendTextMessage,
		"reply_message":         config.dslTelegramReplyTextMessage,
		"send_callback_message": config.dslTelegramSendCallbackMessage,
		"get_bot_updates":       config.dslTelegramGetUpdates,
	}))

	telegramUpdates := L.NewTypeMetatable("telegram_bot_update")
	L.SetGlobal("telegram_bot_update", telegramUpdates)
	L.SetField(telegramUpdates, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"id": config.dslTelegramBotUpdateID,

		"message_id":     config.dslTelegramBotUpdateMessageID,
		"chat_id":        config.dslTelegramBotUpdateMessageChatID,
		"from_id":        config.dslTelegramBotUpdateMessageFromID,
		"from":           config.dslTelegramBotUpdateMessageFromUserName,
		"text":           config.dslTelegramBotUpdateMessageText,
		"reply_text":     config.dslTelegramBotUpdateReplyMessageText,
		"reply_caption":  config.dslTelegramBotUpdateReplyMessageCaption,
		"date":           config.dslTelegramBotUpdateMessageDate,
		"is_bot_command": config.dslTelegramBotUpdateIsBotCommand,

		"is_callback":         config.dslTelegramBotIsCallbackQuery,
		"callback_text":       config.dslTelegramBotCallbackQueryText,
		"callback_from":       config.dslTelegramBotCallbackQueryFromUserName,
		"callback_data":       config.dslTelegramBotCallbackQueryData,
		"callback_message_id": config.dslTelegramBotCallbackQueryMessageID,
	}))

	cache := L.NewTypeMetatable("cache")
	L.SetGlobal("cache", cache)
	L.SetField(cache, "load", L.NewFunction(config.dslNewCache))
	L.SetField(cache, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"set":  config.dslCacheSet,
		"get":  config.dslCacheGet,
		"list": config.dslCacheList,
	}))

	regexp := L.NewTypeMetatable("regexp")
	L.SetGlobal("regexp", regexp)
	L.SetField(regexp, "compile", L.NewFunction(config.dslRegexpCompile))
	L.SetField(regexp, "match", L.NewFunction(config.dslRegexpIsMatch))
	L.SetField(regexp, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"match":           config.dslRegexpMatch,
		"find_all_string": config.dslRegexpFindAllString,
		"find_all":        config.dslRegexpFindAllString,
		"find_string":     config.dslRegexpFindString,
		"find":            config.dslRegexpFindString,
	}))

	postgres := L.NewTypeMetatable("postgres")
	L.SetGlobal("postgres", postgres)
	L.SetField(postgres, "open", L.NewFunction(config.dslNewPgsqlConn))
	L.SetField(postgres, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"close": config.dslPgsqlClose,
		"query": config.dslPgsqlQuery,
	}))

}
