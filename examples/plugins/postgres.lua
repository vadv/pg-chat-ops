package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
settings = require "settings"

local tg_tocken = settings.tg_tocken()
local tg_config = settings.tg_config
local fqdn = settings.fqdn
local notify_chat = settings.notify_chat

tg, err = telegram.new(tg_tocken, tg_config)
if err then error(err) end

chatopsdb, err = cache.load(settings.chat_db)
if err then error(err) end

killdb, err = cache.load(settings.kill_db)
if err then error(err) end

-- открываем "главный" коннект
local main_db, err = postgres.open(settings.connection)
if err then error(err) end

function check_err(err)
  if err then
    main_db:close()
    tg:send_text_message(notify_chat, "Error: `"..err.."` fqdn: "..fqdn, "html")
    error(err)
  end
end

-- устанавливаем лимит на выполнение любого запроса 10s
local _, err = main_db:query("set statement_timeout to '10s'")
check_err(err)

function send_messages_to_all_chats(message)
  local current = chatopsdb:get("chats_id")
  local chats = strings.split(current, ":")
  for _, chat_id in pairs(chats) do
    tg:send_text_message(chat_id, message, "markdown")
  end
end

function send_callback_messages_to_all_chats(message)
  local current = chatopsdb:get("chats_id")
  local chats = strings.split(current, ":")
  local keyboard_inputs = {}
  table.insert(keyboard_inputs, {text= "убить", callback_data = "kill"})
  table.insert(keyboard_inputs, {text= "подождать 30 минут", callback_data = "wait"})
  for _, chat_id in pairs(chats) do
    tg:send_callback_message(chat_id, message, keyboard_inputs, "markdown")
  end
end

function get_pid(query_id)
  local rows, err = main_db:query(" \
    select pid from pg_catalog.pg_stat_activity s where \
      s.state <> 'idle' and \
      md5(s.query || s.query_start::text || s.pid::text) = '"..query_id.."';")
  check_err(err)
  if rows[1] then
    if rows[1][1] then
      return rows[1][1]
    end
  end
  return nil
end

function kill_pid(pid)
  local rows, err = main_db:query("select pg_catalog.pg_cancel_backend("..pid..");")
  check_err(err)
  if rows[1][1] == true then return end
  local rows, err = main_db:query("select pg_catalog.pg_terminate_backend("..pid..");")
  check_err(err)
  if rows[1][1] == true then return end
  return "не смог убить pid: "..pid
end

function check_long_queries()

  local rows, err = main_db:query([[
    select
      s.pid,
      s.query,
      extract(epoch from now() - s.query_start)::int as age,
      s.state,
      s.application_name,
      md5(s.query || s.query_start::text || s.pid::text)
    from
      pg_catalog.pg_stat_activity as s
    where
      not(s.state = 'idle')
      and extract(epoch from now() - s.query_start)::int > 60*60
      and s.query not ilike '%vacuum%'
]])
  check_err(err)
  for _, row in pairs(rows) do
    local pid, query, tt, state = tostring(row[1]), tostring(row[2]), tonumber(row[3]), tostring(row[4])
    local application_name, query_id = tostring(row[5]), tostring(row[6])
    local known_query = killdb:get(query_id)
    if (known_query == nil) then
      local message_template = [[
QueryID:    %s
App name:   `%s`
Запрос:     `%s`
Статус:     `%s`
Время:      `%d сек`

#%s
]]
      local message = string.format(
        message_template,
        query_id,
        application_name,
        query,
        state,
        tt,
        query_id)
      send_callback_messages_to_all_chats(message)
    end

  end
end

function run_killer()
  for query_id, op in pairs(killdb:list()) do
    if op == "kill" then
      local pid = get_pid(query_id)
      if pid then
        local err = kill_pid(pid)
        if err then
          send_messages_to_all_chats("*ОШИБКА* QueryID: `"..query_id.."`: "..err.."\n#"..query_id)
        else
          send_messages_to_all_chats("*УСПЕХ* QueryID: `"..query_id.."`: завершен\n#"..query_id)
          killdb:set(query_id, "miss")
        end
      else
        send_messages_to_all_chats("*ОШИБКА* QueryID: `"..query_id.."`: не найден\n#"..query_id)
        killdb:set(query_id, "miss")
      end
    end
  end
end

local counter = 0
while true do

  chats = chatopsdb:get("chats_id")
  if counter%100 == 0 then if chats then check_long_queries() end end
  if chats then run_killer() end
  counter = counter + 1

  time.sleep(3)
end
