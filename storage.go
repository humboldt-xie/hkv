package main

import (
	"bytes"
	"github.com/syndtr/goleveldb/leveldb"
)

var SlotCount = 16

type Kv interface {
	Set(key []byte, value []byte)
	Get(key []byte) (value []byte, err error)
}

type Mirror interface {
	SetStatus(string) error
	Copy(Item) error
	Sync(Item) error
}

type Item struct {
	Sequence int64
	Key      []byte
	Value    []byte
}

func GetDb() *leveldb.DB {
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		panic(err)
	}
	return db
}

var (
	STATUS_NODE      = "node"
	STATUS_MIGRATING = "migrating"
	STATUS_IMPORTING = "importing"
	STATUS_DELETING  = "deleting"
	STATUS_EMPTY     = "empty"
)

type Dataset struct {
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
	ds.SetStatus(STATUS_MIGRATING)
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
	ds.SetStatus(STATUS_DELETING)
	//ds.Clean()
	err := iter.Error()
	return err
}

type Node struct {
	Id      string
	Dataset map[string]*Dataset
}

func NewNode(Id string) *Node {
	return &Node{Id: Id}
}

var DB = GetDb()

func main() {
	DB.Close()
}
