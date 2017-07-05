package main

import (
	"fmt"
	"sync"
)

type storeValue struct {
	key     []byte
	flags   int
	exptime int
	bytes   int
	data    []byte
}

func (sv *storeValue) toString() string {
	s := fmt.Sprintf("VALUE %s %d %d\n%s\nEND",
		string(sv.key), sv.flags, sv.bytes, string(sv.data))
	return s
}

type inmemoryDB struct {
	mutex  sync.RWMutex
	bucket map[string]*storeValue
}

func (db *inmemoryDB) set(key string, value *storeValue) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.bucket[key] = value

	return nil
}

func (db *inmemoryDB) get(key string) (*storeValue, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	v := db.bucket[key]

	return v, nil
}

func (db *inmemoryDB) add(key string, value *storeValue) (bool, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	v := db.bucket[key]
	if v != nil {
		return false, nil
	}
	db.bucket[key] = value

	return true, nil
}

func (db *inmemoryDB) replace(key string, value *storeValue) (bool, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	v := db.bucket[key]
	if v == nil {
		return false, nil
	}
	db.bucket[key] = value

	return true, nil
}

func (db *inmemoryDB) delete(key string) (bool, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	v := db.bucket[key]
	if v == nil {
		return false, nil
	}
	db.bucket[key] = nil

	return true, nil
}

func (db *inmemoryDB) initialize() error {
	db.bucket = make(map[string]*storeValue)
	return nil
}
