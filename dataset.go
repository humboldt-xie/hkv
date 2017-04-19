package main

import (
	"bytes"
	kv "github.com/humboldt-xie/hkv/proto"
	"log"
	"sync"
)

var (
	STATUS_NODE      = "node"
	STATUS_MIGRATING = "migrating"
	STATUS_IMPORTING = "importing"
	STATUS_DELETING  = "deleting"
	STATUS_EMPTY     = "empty"
)

type Binlog struct {
	Name      []byte
	MaxBinlog int64
}

type Dataset struct {
	mu        sync.Mutex
	Name      string
	Status    string
	Sequence  int64
	MaxBinlog int64
	DbHandle  DbHandle
	mirror    Mirror
}

func (ds *Dataset) RawKey(key []byte) []byte {
	return key[len(ds.Name):]
}

func (ds *Dataset) Key(key []byte) []byte {
	return append([]byte(ds.Name), key...)
}

func (ds *Dataset) SetStatus(status string) error {
	ds.Status = status
	return nil
}

func (ds *Dataset) Copy(data kv.Data) error {
	log.Printf("copy %s", string(ds.Key(data.Key)))
	return ds.DbHandle().Put(ds.Key(data.Key), data.Value, nil)
}

func (ds *Dataset) Sync(data kv.Data) error {
	return ds.DbHandle().Put(ds.Key(data.Key), data.Value, nil)
}

func (ds *Dataset) Set(key []byte, value []byte) error {
	ds.mu.Lock()
	ds.Sequence += 1
	data := kv.Data{Sequence: ds.Sequence, Key: key, Value: value}
	mirror := ds.mirror
	ds.mu.Unlock()
	if mirror != nil {
		mirror.Sync(data)
	}
	return ds.DbHandle().Put(ds.Key(key), value, nil)
}

func (ds *Dataset) Get(key []byte) ([]byte, error) {
	return ds.DbHandle().Get(ds.Key(key), nil)
}
func (ds *Dataset) OnSet(key []byte) bool {
	if len(key) > len(ds.Name) && bytes.Equal([]byte(ds.Name), key[:len(ds.Name)]) {
		return true
	}
	return false
}
func (ds *Dataset) Clean() error {
	iter := ds.DbHandle().NewIterator(nil, nil)
	for ok := iter.Seek([]byte(ds.Name)); ok; iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		if !ds.OnSet(key) {
			break
		}
		ds.DbHandle().Delete(key, nil)
		//value := iter.Value()
		//data := Data{Sequence: ds.Sequence, Key: ds.RawKey(key), Value: value}
		//ds.migrater.Copy(data)
	}
	return nil
}
func (ds *Dataset) SyncTo(sequence int64, mirror Mirror) error {
	return nil
}

func (ds *Dataset) CopyTo(sequence int64, mirror Mirror) error {
	//ds.SetStatus(STATUS_MIGRATING)
	ds.mirror = mirror
	if sequence < ds.Sequence-ds.MaxBinlog || sequence == 0 {
		ds.mirror.SetStatus(STATUS_IMPORTING)
		iter := ds.DbHandle().NewIterator(nil, nil)
		for ok := iter.Seek([]byte(ds.Name)); ok; iter.Next() {
			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			key := iter.Key()
			if !ds.OnSet(key) {
				break
			}
			value := iter.Value()
			data := kv.Data{Sequence: ds.Sequence, Key: ds.RawKey(key), Value: value}
			ds.mirror.Copy(data)
		}
		iter.Release()
		ds.mirror.SetStatus(STATUS_NODE)
		err := iter.Error()
		return err
	}
	return ds.SyncTo(sequence, mirror)
	//ds.SetStatus(STATUS_DELETING)
	//ds.Clean()
}

//func (ds *Dataset) Migrating() error {
//	//ds.SetStatus(STATUS_MIGRATING)
//	ds.migrater.SetStatus(STATUS_IMPORTING)
//	iter := ds.DbHandle().NewIterator(nil, nil)
//	for ok := iter.Seek(ds.Name); ok; iter.Next() {
//		// Remember that the contents of the returned slice should not be modified, and
//		// only valid until the next call to Next.
//		key := iter.Key()
//		if !bytes.Equal(ds.Name, key[:len(ds.Name)]) {
//			break
//		}
//		value := iter.Value()
//		data := Data{Sequence: ds.Sequence, Key: ds.RawKey(key), Value: value}
//		ds.migrater.Copy(data)
//	}
//	iter.Release()
//	ds.migrater.SetStatus(STATUS_NODE)
//	//ds.SetStatus(STATUS_DELETING)
//	//ds.Clean()
//	err := iter.Error()
//	return err
//}
