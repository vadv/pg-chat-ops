package dsl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type dslCaches struct {
	sync.Mutex
	caches map[string]*dslCache
}

var listOfCaches = &dslCaches{caches: make(map[string]*dslCache, 0)}

// cache должен быть общий для всех плагинов
func newDSLCache(filename string) *dslCache {
	listOfCaches.Lock()
	defer listOfCaches.Unlock()
	if result, ok := listOfCaches.caches[filename]; ok {
		return result
	}
	result := &dslCache{filename: filename, List: make(map[string]*dslCacheItem, 0)}
	listOfCaches.caches[filename] = result
	go result.saveRoutine()
	go result.junitorRoutine()
	return result
}

type dslCache struct {
	sync.Mutex
	filename string
	List     map[string]*dslCacheItem `json:"list"`
}
type dslCacheItem struct {
	Value     string `json:"value"`
	ExpiredAt int64  `json:"expired_at"`
}

func (d *dslCacheItem) expired() bool {
	return time.Now().Unix() > d.ExpiredAt
}

func (d *dslCache) load() error {
	d.Lock()
	defer d.Unlock()
	if _, err := os.Stat(d.filename); err != nil {
		return nil
	}
	data, err := ioutil.ReadFile(d.filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}
	d.unsafeClean()
	return nil
}

func (d *dslCache) save() error {
	d.Lock()
	defer d.Unlock()
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	tmpFile := d.filename + ".tmp"
	if err := ioutil.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmpFile, d.filename)
}

func (d *dslCache) unsafeClean() {
	count := 0
	for key, value := range d.List {
		if value.expired() {
			delete(d.List, key)
			count++
		}
	}
	if count > 0 {
		log.Printf("[INFO] delete %d expired items from cache[%s]\n", count, d.filename)
	}
}

func (d *dslCache) count() int {
	d.Lock()
	defer d.Unlock()
	return len(d.List)
}

func (d *dslCache) get(key string) (string, bool) {
	d.Lock()
	defer d.Unlock()
	if value, ok := d.List[key]; ok {
		if !value.expired() {
			return value.Value, true
		}
	}
	return "", false
}

func (d *dslCache) set(key, value string, ttl int64) {
	d.Lock()
	defer d.Unlock()
	d.List[key] = &dslCacheItem{Value: value, ExpiredAt: time.Now().Unix() + ttl}
}

func (d *dslCache) junitorRoutine() {
	for {
		time.Sleep(time.Second * 10)
		d.Lock()
		d.unsafeClean()
		d.Unlock()
	}
}

func (d *dslCache) saveRoutine() {
	for {
		time.Sleep(time.Second * 10)
		if err := d.save(); err != nil {
			log.Printf("[ERROR] save cache[%s]: %s\n", d.filename, err.Error())
		}
		log.Printf("[INFO] saved cache[%s] %d items\n", d.filename, d.count())
	}
}
