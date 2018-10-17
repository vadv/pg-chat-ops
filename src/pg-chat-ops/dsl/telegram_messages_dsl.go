package dsl

import (
	lua "github.com/yuin/gopher-lua"
)

// send_text_message(chat_id, message, <type>)
func (c *dslConfig) dslTelegramSendTextMessage(L *lua.LState) int {
	tg := checkTelegram(L)
	chatID := L.CheckAny(2).String()
	message := L.CheckAny(3).String()
	messageType := "plain"
	if L.GetTop() > 3 {
		messageType = L.CheckAny(4).String()
	}
	if err := tg.sendSimpleTextMessage(chatID, message, messageType); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// reply_message(chat_id, reply_messageID, message, <type>)
func (c *dslConfig) dslTelegramReplyTextMessage(L *lua.LState) int {
	tg := checkTelegram(L)
	chatID := L.CheckAny(2).String()
	messageID := L.CheckAny(3).String()
	message := L.CheckAny(4).String()
	messageType := "plain"
	if L.GetTop() > 4 {
		messageType = L.CheckAny(5).String()
	}
	if err := tg.sendReplyTextMessage(chatID, messageID, message, messageType); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// callback_message(chat_id, message, inline_keyboard, <type>)
func (c *dslConfig) dslTelegramSendCallbackMessage(L *lua.LState) int {
	tg := checkTelegram(L)
	chatID := L.CheckAny(2).String()
	message := L.CheckAny(3).String()
	luaTbl := L.CheckTable(4)
	messageType := "plain"
	if L.GetTop() > 4 {
		messageType = L.CheckAny(5).String()
	}
	goMap := make([]map[string]string, 0)
	luaTbl.ForEach(func(k lua.LValue, v lua.LValue) {
		if tbl, ok := v.(*lua.LTable); ok {
			line := make(map[string]string, 0)
			tbl.ForEach(func(k lua.LValue, v lua.LValue) {
				line[k.String()] = v.String()
			})
			goMap = append(goMap, line)
		}
	})
	if err := tg.sendCallbackMessage(chatID, message, goMap, messageType); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
