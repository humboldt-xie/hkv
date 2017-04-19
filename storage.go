package main

import (
	"fmt"
	kv "github.com/humboldt-xie/hkv/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"time"
)

var SlotCount = 16

type Kv interface {
	Set(key []byte, value []byte)
	Get(key []byte) (value []byte, err error)
}

type Mirror interface {
	SetStatus(string) error
	Copy(kv.Data) error
	Sync(kv.Data) error
}
type DbHandle func() *leveldb.DB

func GetDb() *leveldb.DB {
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		panic(err)
	}
	return db
}

var DB = GetDb()

func GlobalDbHandle() *leveldb.DB {
	return DB
}

func main() {
	s1 := Server{}
	s1.Init("s1")
	log.Printf("init s1")
	go s1.ListenAndServe("127.0.0.1:7001")
	time.Sleep(1 * time.Second)
	s2 := Server{}
	s2.Init("s2")
	log.Printf("init s2")
	go s2.ListenAndServe("127.0.0.1:7002")
	time.Sleep(1 * time.Second)
	log.Printf("init mirror to")
	s2.ImportFrom("127.0.0.1:7001")
	s1.ImportFrom("127.0.0.1:7002")

	log.Printf("set s1")
	s1.Set([]byte("hello"), []byte("world"))
	s2.Set([]byte("hello2"), []byte("world"))

	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 10; i++ {
			s1.Set([]byte(fmt.Sprintf("key-%d", i)), []byte("world"))
		}
	}()

	time.Sleep(2 * time.Second)
	log.Printf("get s2")
	v, err := s2.Get([]byte("hello"))
	log.Printf("s2 hello %s %s", string(v), err)
	v, err = s1.Get([]byte("hello2"))
	log.Printf("s1 hello2 %s %s", string(v), err)
	return
}
