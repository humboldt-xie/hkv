package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	kvproto "github.com/humboldt-xie/hkv/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
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
	mu          sync.Mutex
	Name        string
	Status      string
	MinSequence int64
	Sequence    int64
	MaxBinlog   int64
	dbHandle    DbHandle
	//sync        chan *Data
	mirror Mirror
}

func (ds *Dataset) Init() {
	data, err := ds.LastBinlog()
	if err != nil {
		return
	}
	ds.Sequence = data.Sequence
	first, err := ds.FirstBinlog()
	if err != nil {
		log.Debugf("[%s]get first binlog error:%s", ds.Name, err)
	}
	ds.MinSequence = first.Sequence
	go ds.deleteBinlogBackend()
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

func (ds *Dataset) FirstBinlog() (*kvproto.Data, error) {
	iter := ds.dbHandle().NewIterator(&util.Range{Start: ds.Key("b", []byte{}), Limit: ds.Key("b", []byte(Sequence2str(math.MaxInt64)))}, nil)

	data := &kvproto.Data{}
	if ok := iter.First(); ok {
		err := proto.Unmarshal(iter.Value(), data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (ds *Dataset) LastBinlog() (*kvproto.Data, error) {
	iter := ds.dbHandle().NewIterator(&util.Range{Start: ds.Key("b", []byte{}), Limit: ds.Key("b", []byte(Sequence2str(math.MaxInt64)))}, nil)

	data := &kvproto.Data{}
	if ok := iter.Last(); ok {
		err := proto.Unmarshal(iter.Value(), data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (ds *Dataset) deleteBinlogBackend() {
	for {
		if ds.Sequence-ds.MinSequence > ds.MaxBinlog {
			iter := ds.dbHandle().NewIterator(
				&util.Range{
					Start: ds.Key("b", []byte(Sequence2str(ds.MinSequence))),
					Limit: ds.Key("b", []byte(Sequence2str(ds.Sequence-ds.MaxBinlog))),
				}, nil)

			for ok := iter.First(); ok; iter.Next() {
				data := &kvproto.Data{}
				err := proto.Unmarshal(iter.Value(), data)
				if err == nil {
					ds.MinSequence = data.Sequence
				}
				ds.dbHandle().Delete(iter.Key(), nil)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
func (ds *Dataset) GetBinlog(sequence int64) (*kvproto.Data, error) {
	binlogkey := ds.Key("b", []byte(Sequence2str(sequence)))
	bytes, err := ds.dbHandle().Get(binlogkey, nil)
	if err != nil {
		return nil, err
	}
	data := &kvproto.Data{}
	err = proto.Unmarshal(bytes, data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (ds *Dataset) addBinlog(batch *leveldb.Batch, data *kvproto.Data) ([]byte, error) {
	bytes, err := proto.Marshal(data)
	if err != nil {
		log.Debugf("marshaling error: ", err)
		return nil, err
	}
	binlogkey := ds.Key("b", []byte(Sequence2str(data.Sequence)))
	batch.Put(binlogkey, bytes)

	return bytes, nil
	//return ds.dbHandle().Put(binlogkey, bytes, nil)
}

func (ds *Dataset) SetStatus(status string) error {
	ds.Status = status
	return nil
}

func (ds *Dataset) Copy(data *kvproto.Data) error {
	log.Debugf("copy %s", string(ds.Key("d", data.Key)))
	if data.Sequence <= ds.Sequence {
		//TODO merge data
		return nil
	}
	err := ds.set(data)
	mirror := ds.mirror
	if mirror != nil {
		mirror.Copy(data)
	}
	return err

	//return ds.dbHandle().Put(ds.Key("d", data.Key), data.Value, nil)
}

func (ds *Dataset) Sync(data *kvproto.Data) error {
	if data.Sequence <= ds.Sequence {
		//TODO merge data
		return nil
	}
	err := ds.set(data)
	if err != nil {
		return err
	}
	ds.Sequence = data.Sequence
	mirror := ds.mirror
	if mirror != nil {
		mirror.Sync(data)
	}
	return err
	//return ds.dbHandle().Put(ds.Key("d", data.Key), data.Value, nil)
}

func (ds *Dataset) set(data *kvproto.Data) error {
	db := ds.dbHandle()
	batch := new(leveldb.Batch)
	bytes, err := ds.addBinlog(batch, data)
	if err != nil {
		return err
	}
	batch.Put(ds.Key("d", data.Key), bytes)
	return db.Write(batch, nil)

}

func (ds *Dataset) Set(req *kvproto.SetRequest) (*kvproto.SetReply, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.Sequence += 1
	data := &kvproto.Data{Sequence: ds.Sequence, Key: req.Key, Value: req.Value}
	err := ds.set(data)

	if err != nil {
		return nil, err
	}

	mirror := ds.mirror
	if mirror != nil {
		mirror.Sync(data)
	}

	return &kvproto.SetReply{Sequence: data.Sequence}, nil
}
func (ds *Dataset) ConfigSet(key string, value string) error {
	err := ds.dbHandle().Put(ds.Key("c", []byte(key)), []byte(value), nil)
	return err
}
func (ds *Dataset) ConfigGet(key string) string {
	bytes, err := ds.dbHandle().Get(ds.Key("c", []byte(key)), nil)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (ds *Dataset) GetData(key []byte) (*kvproto.Data, error) {
	bytes, err := ds.dbHandle().Get(ds.Key("d", key), nil)
	if err != nil {
		return nil, err
	}
	data := &kvproto.Data{}
	err = proto.Unmarshal(bytes, data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (ds *Dataset) Get(req *kvproto.GetRequest) (*kvproto.GetReply, error) {
	data, err := ds.GetData(req.Key)
	if err != nil {
		return nil, err
	}
	return &kvproto.GetReply{Data: data}, nil
}

func (ds *Dataset) OnDataset(key []byte) bool {
	if len(key) > len(ds.Name) && bytes.Equal([]byte(ds.Name), key[:len(ds.Name)]) {
		return true
	}
	return false
}
func (ds *Dataset) Clean() error {
	iter := ds.dbHandle().NewIterator(nil, nil)
	for ok := iter.Seek([]byte(ds.Name)); ok; iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		if !ds.OnDataset(key) {
			break
		}
		ds.dbHandle().Delete(key, nil)
	}
	return nil
}

// sync binlog to
func (ds *Dataset) SyncTo(sequence int64, mirror Mirror) error {
	log.Debugf("[%s]SyncTo [%d-%d]", ds.Name, sequence, ds.Sequence)
	for ; sequence <= ds.Sequence; sequence++ {
		data, err := ds.GetBinlog(sequence)
		if err != nil {
			log.Debugf("[%s] get binlog[%d] error:%s", ds.Name, sequence, err)
			//TODO fix read binlog error
			continue
		}
		mirror.Sync(data)
	}
	return nil
}

func (ds *Dataset) CopyTo(sequence int64, mirror Mirror) error {
	ds.mirror = mirror
	log.Debugf("[%s]CopyTo [%d-%d]", ds.Name, sequence, ds.Sequence)
	if sequence < ds.MinSequence || sequence == 0 {
		ds.mirror.SetStatus(STATUS_IMPORTING)
		iter := ds.dbHandle().NewIterator(nil, nil)
		for ok := iter.Seek(ds.Key("d", []byte{})); ok; iter.Next() {

			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			key := iter.Key()
			if !ds.OnDataset(key) {
				break
			}
			value := iter.Value()
			data := &kvproto.Data{}
			err := proto.Unmarshal(value, data)
			if err != nil {
				log.Debugf("data error:%s %s", string(key), err)
			}

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
