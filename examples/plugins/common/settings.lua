local settings = {}

settings.fqdn = strings.trim(ioutil.readfile("/proc/sys/kernel/hostname"), "\n")
settings.tg_config = {
--  proxy      = "http://",
  ignore_ssl = true,
}

settings.notify_chat = XXXX 
settings.chat_db = "chatops_db.json"
settings.kill_db = "kill_db.json"

function settings.tg_tocken()
  if string.match(settings.fqdn, "apidb01") then return "73XXX:XXXX" end
end

-- для работы данного плагина необходимо
-- 1. создать пользователя root `create user root with superuser;`
-- 2. дать разрешение в pg_hba: `local all root peer` не надо притворяться что unix-root не имеет полный доступ в базу
local connection = {
  host     = '/tmp',
  user     = 'root',
  database = 'postgres'
}
if goos.stat('/var/run/postgresql/.s.PGSQL.5432') then connection.host = '/var/run/postgresql' end

settings.connection = connection

return settings
