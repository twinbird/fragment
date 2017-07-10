package main

import (
	"fmt"
	"sync"
	"time"
)

type storeValue struct {
	key        []byte
	flags      int
	exptime    int
	bytes      int
	data       []byte
	storedTime time.Time
}

func (sv *storeValue) toString() string {
	s := fmt.Sprintf("VALUE %s %d %d\r\n%s\r\nEND\r\n",
		string(sv.key), sv.flags, sv.bytes, string(sv.data))
	return s
}

type inmemoryDB struct {
	mutex  sync.RWMutex
	bucket map[string]*storeValue
}

func isExpired(v *storeValue) bool {
	if v.exptime == 0 {
		return false
	}
	exptime := time.Duration(v.exptime)
	t := v.storedTime.Add(exptime * time.Second)
	if t.After(time.Now()) == false {
		return true
	} else {
		return false
	}
}

func (db *inmemoryDB) set(key []byte, value *storeValue) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	value.storedTime = time.Now()
	db.bucket[string(key)] = value

	return nil
}

func (db *inmemoryDB) get(key []byte) (*storeValue, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	v := db.bucket[string(key)]
	if v == nil {
		return nil, nil
	}

	if isExpired(v) == true {
		return nil, nil
	}

	return v, nil
}

func (db *inmemoryDB) add(key []byte, value *storeValue) (bool, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	value.storedTime = time.Now()
	v := db.bucket[string(key)]
	if v != nil {
		if isExpired(v) == false {
			return false, nil
		}
	}

	db.bucket[string(key)] = value

	return true, nil
}

func (db *inmemoryDB) replace(key []byte, value *storeValue) (bool, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	value.storedTime = time.Now()
	v := db.bucket[string(key)]
	if v == nil {
		return false, nil
	}

	if isExpired(v) == true {
		return false, nil
	}

	db.bucket[string(key)] = value

	return true, nil
}

func (db *inmemoryDB) delete(key []byte) (bool, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	v := db.bucket[string(key)]
	if v == nil {
		return false, nil
	}

	if isExpired(v) == true {
		return false, nil
	}

	db.bucket[string(key)] = nil

	return true, nil
}

func (db *inmemoryDB) initialize() error {
	db.bucket = make(map[string]*storeValue)
	return nil
}
