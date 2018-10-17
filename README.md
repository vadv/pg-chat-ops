# intro

толстый golang-бинарь, который запускает embeded lua-скрипты, для управления и информации pg 

В examples лежит пример на котором можно построить telegram-бота, для того чтобы управлять длинными транзакциями

![Пример](https://i.imgur.com/FZfsPgL.png)

# lua-dsl

## postgres

```
db, err = postgres.open({database="xxx", host="127.0.0.1", user="xxx", password="xxx"}) открыть коннект
rows, err, column_count, row_count = db:query() выполнить запрос
db:close() закрыть коннект
```

## telegram

```
config = {
    -- proxy = = "http://user:password@proxy",
    -- ignore_ssl = true/false,
}

tg, err = telegram.new(tocken, config)
tg:send_text_message(chat_id, message, <"markdown"|"html">)
tg:reply_message(chat_id, message_id, message, <"markdown"|"html">)
tg:send_zabbix_graph(zbx_trigger, chat_id, message, <"markdown"|"html">)
tg:send_callback_message(chat_id, message, { {text = "", callback_data = ""},{text = "", callback_data = ""}}, <"markdown"|"html">)

tg:get_bot_updates() -- table из user-data telegram-bot-updates (назовем tg_bot_update)

tg_bot_update:id()
tg_bot_update:message_id()
tg_bot_update:chat_id()
tg_bot_update:from_id()
tg_bot_update:from() -- username
tg_bot_update:text()
tg_bot_update:reply_text() -- message reply
tg_bot_update:reply_caption() -- message reply caption
tg_bot_update:date() -- unixts
tg_bot_update:is_bot_command true/false

tg_bot_update:is_callback() -- ответ на callback
tg_bot_update:callback_text() -- text сообщения на который был callback
tg_bot_update:callback_from() -- имя пользователя который послал callback
tg_bot_update:callback_data() -- данные который отослал callback
tg_bot_update:callback_message_id() -- message id запроса который и породил callback
```

## cache (для обмена между плагинами и сохранением состояния)

```
c, err = cache.load(filename)

c:set(key, value, <ttl>)
c:get(key) -- nil/value(string)
c:list() -- { k1="v1", k2="v2"}
```
