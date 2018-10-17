package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
settings = require "settings"

local tg_tocken = settings.tg_tocken()
local tg_config = settings.tg_config
local fqdn = settings.fqdn
local notify_chat = settings.notify_chat

tg, err = telegram.new(tg_tocken, tg_config)
if err then error(err) end

err = tg:send_text_message(notify_chat, "i'm online, fqdn: "..fqdn, "html")
if err then error(err) end

chatopsdb, err = cache.load(settings.chat_db)
if err then error(err) end

killdb, err = cache.load(settings.kill_db)
if err then error(err) end

function process_command(upd)
  local text = upd:text()
  log.info("process command: "..text)
  if strings.has_prefix(text, "/help") then help(upd) end
  if strings.has_prefix(text, "/chat_id") then chat_id(upd) end
  if strings.has_prefix(text, "/start") then start(upd) end
  if strings.has_prefix(text, "/stop") then stop(upd) end
end

function help(upd)
  log.info("start telegram bot command: help")
  local help_message = [[
/help       - эта помощь
/chat_id    - сообщить chat_id

/start - запустить бота в этом чате
/stop  - остановить бота в этом чате

  ]]
  err = tg:send_text_message(upd:chat_id(), help_message)
  if err then error(err) end
end

function chat_id(upd)
  err = tg:reply_message(upd:chat_id(), upd:message_id(), upd:chat_id())
  if err then error(err) end
end

function start(upd)
  local chat_id = tostring(upd:chat_id())
  local current = chatopsdb:get("chats_id")
  if current then
    local found_in_db = false
    for _, v in pairs(strings.split(current, ":")) do if v == chat_id then found_in_db = true end end
    if not found_in_db then
      current = current..":"..chat_id
    else
      tg:reply_message(upd:chat_id(), upd:message_id(), "такой чат уже зарегистрирован")
      return
    end
  else
    current = chat_id
  end
  chatopsdb:set("chats_id", current, 10000000000)
  log.info("set chats_id: "..current)
  tg:reply_message(upd:chat_id(), upd:message_id(), "готово, для ChatID: "..chat_id)
end

function stop(upd)
  local chat_id = tostring(upd:chat_id())
  local current = chatopsdb:get("chats_id")
  if current then
    local found_in_db = false
    for _, v in pairs(strings.split(current, ":")) do if v == chat_id then found_in_db = true end end
    if found_in_db then
      -- удаляем
      local new_current = ""
      for _, v in pairs(strings.split(current, ":")) do
        if not(v == chat_id) then
          if not(v == "") then
            new_current = new_current .. ":" .. v
          end
        end
      end
      chatopsdb:set("chats_id", new_current, 10000000000)
      log.info("set chats_id: "..new_current)
      tg:reply_message(upd:chat_id(), upd:message_id(), "готово, со следующей пачки")
    else
      tg:reply_message(upd:chat_id(), upd:message_id(), "такой чат не зарегистрирован")
    end
  else
    tg:reply_message(upd:chat_id(), upd:message_id(), "такой чат не зарегистрирован: база пуста")
  end
end

function process_callback(upd)
  local data = upd:callback_data()
  if data == "kill" then set_kill_query(upd) end
  if data == "wait" then set_wait_query(upd) end
end

function send_messages_to_all_chats(message)
  local current = chatopsdb:get("chats_id")
  local chats = strings.split(current, ":")
  for _, chat_id in pairs(chats) do
    print(chat_id, message)
    tg:send_text_message(chat_id, message, "markdown")
  end
end

function get_query_id(text)
  return string.match(text, "QueryID:%s+(%S+)")
end

function set_kill_query(upd)
  local text = upd:callback_text()
  local query_id = get_query_id(text)
  if query_id then
    killdb:set(query_id, "kill", 300)
    send_messages_to_all_chats("QueryID: "..query_id.." будет убит (от @"..upd:callback_from()..")")
  else
    send_messages_to_all_chats("QueryID не найден в "..text)
  end
end

function set_wait_query(upd)
  local text = upd:callback_text()
  local query_id = get_query_id(text)
  if query_id then
    killdb:set(query_id, "wait", 30*60)
    send_messages_to_all_chats("QueryID: "..query_id.." отложено на 30 минут (от @"..upd:callback_from()..")")
  else
    send_messages_to_all_chats("QueryID не найден в \n"..text)
  end
end

function run_update()

  local updates, err = tg:get_bot_updates()
  if err then error(err) end

  for _, upd in pairs(updates) do
    if upd:is_bot_command() then process_command(upd) end
    if upd:is_callback() then process_callback(upd) end
  end

end

while true do
  run_update()
  time.sleep(1)
end
