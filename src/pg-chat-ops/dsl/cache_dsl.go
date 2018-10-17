package dsl

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// получение cache из lua-state
func checkCache(L *lua.LState) *dslCache {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*dslCache); ok {
		return v
	}
	L.ArgError(1, "cache expected")
	return nil
}

func (c *dslConfig) dslNewCache(L *lua.LState) int {
	filename := L.CheckString(1)
	d := newDSLCache(filename)
	if err := d.load(); err != nil {
		log.Printf("[ERROR] load cache[%s] error: %s\n", d.filename, err.Error())
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	log.Printf("[INFO] loaded cache[%s] %d items\n", d.filename, d.count())
	ud := L.NewUserData()
	ud.Value = d
	L.SetMetatable(ud, L.GetTypeMetatable("cache"))
	L.Push(ud)
	log.Printf("[INFO] New cache[%s]\n", d.filename)
	return 1
}

func (c *dslConfig) dslCacheSet(L *lua.LState) int {
	d := checkCache(L)
	key := L.CheckAny(2).String()
	value := L.CheckAny(3).String()
	ttl := int64(60)
	if L.GetTop() > 3 {
		ttl = int64(L.CheckNumber(4))
	}
	d.set(key, value, ttl)
	return 0
}

func (c *dslConfig) dslCacheGet(L *lua.LState) int {
	d := checkCache(L)
	key := L.CheckAny(2).String()
	val, ok := d.get(key)
	if ok {
		L.Push(lua.LString(val))
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

func (c *dslConfig) dslCacheList(L *lua.LState) int {

	d := checkCache(L)
	d.Lock()
	defer d.Unlock()

	result := L.NewTable()
	for k, v := range d.List {
		result.RawSetString(k, lua.LString(v.Value))
	}
	L.Push(result)
	return 1
}
