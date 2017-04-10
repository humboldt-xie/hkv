package main

import (
	"bytes"
	"sync"
)

var (
	STATUS_NODE      = "node"
	STATUS_MIGRATING = "migrating"
	STATUS_IMPORTING = "importing"
	STATUS_DELETING  = "deleting"
	STATUS_EMPTY     = "empty"
)

type Dataset struct {
	mu       sync.Mutex
	Name     []byte
	Status   string
	Sequence int64
	migrater Mirror
}

func (ds *Dataset) RawKey(key []byte) []byte {
	return key[len(ds.Name):]
}

func (ds *Dataset) Key(key []byte) []byte {
	return append(ds.Name, key...)
}

func (ds *Dataset) SetStatus(status string) error {
	ds.Status = status
	return nil
}

func (ds *Dataset) Copy(item Item) error {
	return DB.Put(ds.Key(item.Key), item.Value, nil)
}

func (ds *Dataset) Sync(item Item) error {
	return DB.Put(ds.Key(item.Key), item.Value, nil)
}

func (ds *Dataset) Set(key []byte, value []byte) error {
	ds.Sequence += 1
	item := Item{Sequence: ds.Sequence, Key: key, Value: value}
	if ds.migrater != nil {
		ds.migrater.Sync(item)
	}
	return DB.Put(ds.Key(key), value, nil)
}

func (ds *Dataset) Get(key []byte) ([]byte, error) {
	return DB.Get(ds.Key(key), nil)
}
func (ds *Dataset) OnSet(key []byte) bool {
	if !bytes.Equal(ds.Name, key[:len(ds.Name)]) {
		return false
	}
	return true
}
func (ds *Dataset) Clean() error {
	iter := DB.NewIterator(nil, nil)
	for ok := iter.Seek(ds.Name); ok; iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		if !ds.OnSet(key) {
			break
		}
		DB.Delete(key, nil)
		//value := iter.Value()
		//item := Item{Sequence: ds.Sequence, Key: ds.RawKey(key), Value: value}
		//ds.migrater.Copy(item)
	}
	return nil
}

func (ds *Dataset) Migrating() error {
	//ds.SetStatus(STATUS_MIGRATING)
	ds.migrater.SetStatus(STATUS_IMPORTING)
	iter := DB.NewIterator(nil, nil)
	for ok := iter.Seek(ds.Name); ok; iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		if !bytes.Equal(ds.Name, key[:len(ds.Name)]) {
			break
		}
		value := iter.Value()
		item := Item{Sequence: ds.Sequence, Key: ds.RawKey(key), Value: value}
		ds.migrater.Copy(item)
	}
	iter.Release()
	ds.migrater.SetStatus(STATUS_NODE)
	//ds.SetStatus(STATUS_DELETING)
	//ds.Clean()
	err := iter.Error()
	return err
}
