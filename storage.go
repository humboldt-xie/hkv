package main

import (
	kvproto "github.com/humboldt-xie/hkv/proto"
	"github.com/syndtr/goleveldb/leveldb"

	"fmt"
	"log"
	"time"
)

var SlotCount = 16

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
	s1.AddDataset("1-")
	s1.AddDataset("2-")
	log.Printf("init s1")
	go s1.ListenAndServe("127.0.0.1:7001")
	time.Sleep(1 * time.Second)
	s2 := Server{}
	s2.Init("s2")
	s2.AddDataset("1-")
	s2.AddDataset("2-")
	log.Printf("init s2")
	go s2.ListenAndServe("127.0.0.1:7002")
	time.Sleep(1 * time.Second)
	log.Printf("init mirror to")
	s2.ImportFrom("127.0.0.1:7001")
	s1.ImportFrom("127.0.0.1:7002")

	log.Printf("set s1")
	s1.Set(&kvproto.SetRequest{Dataset: "1-", Key: []byte("hello"), Value: []byte("world")})
	s2.Set(&kvproto.SetRequest{Dataset: "1-", Key: []byte("hello2"), Value: []byte("world")})

	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 10; i++ {
			//s1.Set([]byte(fmt.Sprintf("key-%d", i)), []byte("world"))
			key := []byte(fmt.Sprintf("key-%d", i))
			s1.Set(&kvproto.SetRequest{Dataset: "1-", Key: key, Value: []byte("world")})
		}
		for i := 0; i < 10; i++ {
			//s1.Set([]byte(fmt.Sprintf("key-%d", i)), []byte("world"))
			key := []byte(fmt.Sprintf("key-%d", i))
			s1.Set(&kvproto.SetRequest{Dataset: "2-", Key: key, Value: []byte("world")})
		}
	}()

	time.Sleep(2 * time.Second)
	log.Printf("get s2")
	v, err := s2.Get(&kvproto.GetRequest{Dataset: "1-", Key: []byte("hello")})
	log.Printf("s2 hello %#v %s", v, err)
	v, err = s1.Get(&kvproto.GetRequest{Dataset: "1-", Key: []byte("hello")})
	log.Printf("s1 hello2 %#v %s", v, err)
	log.Printf("s1 info")
	log.Printf("%s", s1.Info())
	log.Printf("s2 info")
	log.Printf("%s", s2.Info())
	time.Sleep(1000 * time.Second)
	return
}
