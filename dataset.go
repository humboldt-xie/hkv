package main

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	kv "github.com/humboldt-xie/hkv/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
)

var (
	STATUS_NODE      = "node"
	STATUS_MIGRATING = "migrating"
	STATUS_IMPORTING = "importing"
	STATUS_DELETING  = "deleting"
	STATUS_EMPTY     = "empty"
)

func Sequence2str(seq int64) string {
	return fmt.Sprintf("%16x", seq)
}
func Sequence2Int(seq string) int64 {
	seq = strings.Trim(seq, " ")
	useq, _ := strconv.ParseInt(seq, 16, 64)
	return useq
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

func (ds *Dataset) Init() {
	data, err := ds.LastBinlog()
	if err != nil {
		return
	}
	ds.Sequence = data.Sequence
}

func (ds *Dataset) RawKey(t string, key []byte) []byte {
	if len(key) < len(ds.Name+t) {
		return nil
	}
	return key[len(ds.Name+t):]
}

func (ds *Dataset) Key(t string, key []byte) []byte {
	return append([]byte(ds.Name+t), key...)
}

func (ds *Dataset) LastBinlog() (*kv.Data, error) {
	iter := ds.DbHandle().NewIterator(&util.Range{Start: ds.Key("b", []byte{}), Limit: ds.Key("b", []byte(Sequence2str(math.MaxInt64)))}, nil)

	data := &kv.Data{}
	if ok := iter.Last(); ok {
		err := proto.Unmarshal(iter.Value(), data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (ds *Dataset) AddBinlog(batch *leveldb.Batch, data *kv.Data) ([]byte, error) {
	bytes, err := proto.Marshal(data)
	if err != nil {
		log.Printf("marshaling error: ", err)
		return nil, err
	}
	binlogkey := ds.Key("b", []byte(Sequence2str(data.Sequence)))
	batch.Put(binlogkey, bytes)

	return bytes, nil
	//return ds.DbHandle().Put(binlogkey, bytes, nil)
}

func (ds *Dataset) SetStatus(status string) error {
	ds.Status = status
	return nil
}

func (ds *Dataset) Copy(data kv.Data) error {
	log.Printf("copy %s", string(ds.Key("d", data.Key)))
	if data.Sequence <= ds.Sequence {
		//TODO merge data
		return nil
	}
	err := ds.set(&data)
	mirror := ds.mirror
	if mirror != nil {
		mirror.Copy(data)
	}
	return err

	//return ds.DbHandle().Put(ds.Key("d", data.Key), data.Value, nil)
}

func (ds *Dataset) Sync(data kv.Data) error {
	if data.Sequence <= ds.Sequence {
		//TODO merge data
		return nil
	}
	err := ds.set(&data)
	if err != nil {
		return err
	}
	ds.Sequence = data.Sequence
	mirror := ds.mirror
	if mirror != nil {
		mirror.Sync(data)
	}
	return err
	//return ds.DbHandle().Put(ds.Key("d", data.Key), data.Value, nil)
}

func (ds *Dataset) set(data *kv.Data) error {
	db := ds.DbHandle()
	batch := new(leveldb.Batch)
	bytes, err := ds.AddBinlog(batch, data)
	if err != nil {
		return err
	}
	batch.Put(ds.Key("d", data.Key), bytes)
	return db.Write(batch, nil)

}

func (ds *Dataset) Set(key []byte, value []byte) error {
	ds.mu.Lock()
	ds.Sequence += 1
	data := &kv.Data{Sequence: ds.Sequence, Key: key, Value: value}
	err := ds.set(data)
	ds.mu.Unlock()

	mirror := ds.mirror
	if mirror != nil {
		mirror.Sync(*data)
	}

	return err
}

func (ds *Dataset) GetData(key []byte) (*kv.Data, error) {
	bytes, err := ds.DbHandle().Get(ds.Key("d", key), nil)
	if err != nil {
		return nil, err
	}
	data := &kv.Data{}
	err = proto.Unmarshal(bytes, data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (ds *Dataset) Get(key []byte) ([]byte, error) {
	data, err := ds.GetData(key)
	if err != nil {
		return nil, err
	}
	return data.Value, nil
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
		for ok := iter.Seek(ds.Key("d", []byte{})); ok; iter.Next() {

			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			key := iter.Key()
			if !ds.OnSet(key) {
				break
			}
			value := iter.Value()
			data := &kv.Data{}
			err := proto.Unmarshal(value, data)
			if err != nil {
				log.Printf("data error:%s %s", string(key), err)
			}

			ds.mirror.Copy(*data)
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
