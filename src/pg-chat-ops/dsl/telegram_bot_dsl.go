package dsl

import (
	lua "github.com/yuin/gopher-lua"
)

// получение bot update из lua-state
func checkTelegramBotUpdate(L *lua.LState, top int) *dslTelegramBotUpdate {
	ud := L.CheckUserData(top)
	if v, ok := ud.Value.(*dslTelegramBotUpdate); ok {
		return v
	}
	L.ArgError(top, "telegram bot update expected")
	return nil
}

func (c *dslConfig) dslTelegramBotUpdateID(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	L.Push(lua.LNumber(upd.UpdateID))
	return 1
}

func (c *dslConfig) dslTelegramBotIsCallbackQuery(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	has := upd.CallbackQuery != nil
	L.Push(lua.LBool(has))
	return 1
}

func (c *dslConfig) dslTelegramBotCallbackQueryID(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.CallbackQuery != nil {
		L.Push(lua.LString(upd.CallbackQuery.ID))
		return 1
	}
	return 0
}

func (c *dslConfig) dslTelegramBotCallbackChatID(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.CallbackQuery != nil {
		if upd.CallbackQuery.Message != nil {
			if upd.CallbackQuery.Message.Chat != nil {
				L.Push(lua.LString(upd.CallbackQuery.Message.Chat.ID))
				return 1
			}
		}
	}
	return 0
}

func (c *dslConfig) dslTelegramBotCallbackQueryData(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.CallbackQuery != nil {
		L.Push(lua.LString(upd.CallbackQuery.Data))
		return 1
	}
	return 0
}

func (c *dslConfig) dslTelegramBotCallbackQueryFromUserName(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.CallbackQuery != nil {
		if upd.CallbackQuery.From != nil {
			L.Push(lua.LString(upd.CallbackQuery.From.UserName))
			return 1
		}
	}
	return 0
}

func (c *dslConfig) dslTelegramBotCallbackQueryMessageID(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.CallbackQuery != nil {
		if upd.CallbackQuery.Message != nil {
			L.Push(lua.LNumber(upd.CallbackQuery.Message.MessageID))
			return 1
		}
	}
	return 0
}

func (c *dslConfig) dslTelegramBotCallbackQueryText(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.CallbackQuery != nil {
		if upd.CallbackQuery.Message != nil {
			L.Push(lua.LString(upd.CallbackQuery.Message.Text))
			return 1
		}
	}
	return 0
}

func (c *dslConfig) dslTelegramBotUpdateMessageID(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.Message != nil {
		L.Push(lua.LNumber(upd.Message.MessageID))
		return 1
	}
	return 0
}

func (c *dslConfig) dslTelegramBotUpdateMessageChatID(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	//upd.Message.From.ID
	if upd.Message != nil {
		if upd.Message.Chat != nil {
			L.Push(lua.LNumber(upd.Message.Chat.ID))
			return 1
		}
	}
	return 0
}

func (c *dslConfig) dslTelegramBotUpdateMessageFromID(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	//upd.Message.From.ID
	if upd.Message != nil {
		if upd.Message.From != nil {
			L.Push(lua.LNumber(upd.Message.From.ID))
			return 1
		}
	}
	return 0
}

func (c *dslConfig) dslTelegramBotUpdateMessageFromUserName(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	//upd.Message.From.ID
	if upd.Message != nil {
		if upd.Message.From != nil {
			L.Push(lua.LString(upd.Message.From.UserName))
			return 1
		}
	}
	return 0
}

func (c *dslConfig) dslTelegramBotUpdateMessageText(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.Message != nil {
		L.Push(lua.LString(upd.Message.Text))
		return 1
	}
	return 0
}

func (c *dslConfig) dslTelegramBotUpdateReplyMessageText(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.Message != nil {
		if upd.Message.ReplyToMessage != nil {
			L.Push(lua.LString(upd.Message.ReplyToMessage.Text))
			return 1
		}
	}
	return 0
}

func (c *dslConfig) dslTelegramBotUpdateReplyMessageCaption(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.Message != nil {
		if upd.Message.ReplyToMessage != nil {
			L.Push(lua.LString(upd.Message.ReplyToMessage.Caption))
			return 1
		}
	}
	return 0
}

func (c *dslConfig) dslTelegramBotUpdateMessageDate(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.Message != nil {
		L.Push(lua.LNumber(upd.Message.Date))
		return 1
	}
	return 0
}

func (c *dslConfig) dslTelegramBotUpdateIsBotCommand(L *lua.LState) int {
	upd := checkTelegramBotUpdate(L, 1)
	if upd.Message != nil {
		if upd.Message.Entities != nil {
			entities := *upd.Message.Entities
			for _, ent := range entities {
				if ent.Type == `bot_command` {
					L.Push(lua.LBool(true))
					return 1
				}
			}
		}
	}
	L.Push(lua.LBool(false))
	return 1
}

func (c *dslConfig) dslTelegramGetUpdates(L *lua.LState) int {
	tg := checkTelegram(L)
	updates, err := tg.getUpdates()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	for _, upd := range updates {
		ud := L.NewUserData()
		ud.Value = upd
		L.SetMetatable(ud, L.GetTypeMetatable("telegram_bot_update"))
		result.Append(ud)
	}
	L.Push(result)
	return 1
}
