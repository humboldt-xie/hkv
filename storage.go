package main

import (
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
