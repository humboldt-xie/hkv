package main

import (
	"testing"
)

type printer struct {
	Name string
	t    *testing.T
}

func (ds *printer) SetStatus(status string) error {
	ds.t.Log(ds.Name, ":set status", status)
	return nil
}

func (ds *printer) Copy(item Item) error {
	ds.t.Log(ds.Name, ":copy", string(item.Key), string(item.Value))
	return nil
}

func (ds *printer) Sync(item Item) error {
	ds.t.Log(ds.Name, ":sync", string(item.Key), string(item.Value))
	return nil
}

type LocalDatasetMigrater struct {
	t       *testing.T
	other   *Dataset
	lastkey string
}

func (ds *LocalDatasetMigrater) SetStatus(status string) error {
	ds.t.Log("set status", status)
	return ds.other.SetStatus(status)
	//return nil
}

func (ds *LocalDatasetMigrater) Copy(item Item) error {
	ds.lastkey = string(item.Key)
	ds.t.Log("copy[", ds.lastkey, "]", string(item.Key), string(item.Value))
	return ds.other.Copy(item)
	//return nil //DB.Put(item.Key, item.Value)
}

func (ds *LocalDatasetMigrater) Sync(item Item) error {
	ds.t.Log("sync[", ds.lastkey, "]", string(item.Key), string(item.Value))
	return ds.other.Sync(item)
	//return nil //DB.Put(item.Key, item.Value)
}

func TestSlot(t *testing.T) {
	other := Dataset{Name: []byte("other-"), Status: "node", migrater: &printer{Name: "other", t: t}}
	other.Clean()
	other.Migrating()

	ds := Dataset{Name: []byte("test-"), Status: "node", migrater: &LocalDatasetMigrater{t: t, other: &other}}
	ds.Set([]byte("hello"), []byte("world"))
	ds.Set([]byte("hello2"), []byte("world2"))
	ds.Migrating()
	ds.Set([]byte("hello1"), []byte("world1"))
	ds.Set([]byte("hello"), []byte("world-second"))

	//other.migrater = &printer{t: t}
	other.Migrating()
	t.Fatal("end")
}
