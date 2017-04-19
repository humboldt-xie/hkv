package main

import (
	hkv "github.com/humboldt-xie/hkv/proto"
	"testing"
)

type LocalDatasetMigrater struct {
	Name    string
	t       *testing.T
	other   *Dataset
	lastkey string
}

func (ds *LocalDatasetMigrater) SetStatus(status string) error {
	ds.t.Log("set status", status)
	if ds.other == nil {
		return nil
	}
	return ds.other.SetStatus(status)
	//return nil
}

func (ds *LocalDatasetMigrater) Copy(data hkv.Data) error {
	ds.lastkey = string(data.Key)
	ds.t.Log("copy[", ds.lastkey, "](", data.Sequence, ")", string(data.Key), string(data.Value))
	if ds.other == nil {
		return nil
	}
	return ds.other.Copy(data)
	//return nil //DB.Put(data.Key, data.Value)
}

func (ds *LocalDatasetMigrater) Sync(data hkv.Data) error {
	ds.t.Log("sync[", ds.lastkey, "](", data.Sequence, ")", string(data.Key), string(data.Value))
	if ds.other == nil {
		return nil
	}
	return ds.other.Sync(data)
	//return nil //DB.Put(data.Key, data.Value)
}

func TestSlot(t *testing.T) {
	p := &LocalDatasetMigrater{Name: "other", t: t}
	other := Dataset{DbHandle: GlobalDbHandle, Name: "other-", Status: "node"}
	other.Clean()
	other.CopyTo(0, p)

	ds := Dataset{DbHandle: GlobalDbHandle, Name: "test-", Status: "node"}
	ds.Init()
	ds.Set([]byte("hello"), []byte("world"))
	ds.Set([]byte("hello2"), []byte("world2"))
	ds.CopyTo(0, &LocalDatasetMigrater{t: t, other: &other})
	ds.Set([]byte("hello1"), []byte("world1"))
	ds.Set([]byte("hello"), []byte("world-second"))

	//other.migrater = &printer{t: t}
	other.CopyTo(0, p)

	data, err := ds.LastBinlog()
	t.Log("last binlog : %#v %s", data, err)
	t.Fatal("end")
}
