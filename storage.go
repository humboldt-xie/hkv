package main

import (
	log "github.com/Sirupsen/logrus"
	kvproto "github.com/humboldt-xie/hkv/proto"
	"github.com/syndtr/goleveldb/leveldb"

	"fmt"
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
	log.SetLevel(log.DebugLevel)

	servers := []*Server{}
	count := 2

	for i := 0; i < count; i++ {
		s := &Server{}
		s.Init(fmt.Sprintf("/data/s%d", i))
		s.AddDataset("1-")
		s.AddDataset("2-")
		log.Debugf("init s[%d]", i)
		go s.ListenAndServe("127.0.0.1:" + fmt.Sprintf("700%d", i))
		servers = append(servers, s)
		for j := 0; j < count; j++ {
			if i != j {
				addr := "127.0.0.1:" + fmt.Sprintf("700%d", i)
				servers[i].ImportFrom("1-", addr)
				servers[i].ImportFrom("2-", addr)
			}
		}
	}

	s1 := servers[0]
	s2 := servers[1]
	//time.Sleep(1 * time.Second)
	//s2 := Server{}
	//s2.Init("s2")
	//s2.AddDataset("1-")
	//s2.AddDataset("2-")
	//log.Debugf("init s2")
	//go s2.ListenAndServe("127.0.0.1:7002")
	//time.Sleep(1 * time.Second)
	//log.Debugf("init mirror to")
	//s2.ImportFrom("1-", "127.0.0.1:7001")
	//s2.ImportFrom("2-", "127.0.0.1:7001")
	//s1.ImportFrom("2-", "127.0.0.1:7002")
	//s1.ImportFrom("1-", "127.0.0.1:7002")

	//log.Debugf("set s1")
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
	go func() {
		for i := 0; ; i++ {
			//s1.Set([]byte(fmt.Sprintf("key-%d", i)), []byte("world"))
			key := []byte(fmt.Sprintf("sleep-%d", i))
			s1.Set(&kvproto.SetRequest{Dataset: "1-", Key: key, Value: []byte("world")})
			time.Sleep(1 * time.Second)
		}

	}()
	time.Sleep(2 * time.Second)
	log.Debugf("get s2")
	v, err := s2.Get(&kvproto.GetRequest{Dataset: "1-", Key: []byte("hello")})
	log.Debugf("s2 hello %#v %s", v, err)
	v, err = s1.Get(&kvproto.GetRequest{Dataset: "1-", Key: []byte("hello")})
	log.Debugf("s1 hello2 %#v %s", v, err)
	for i := 0; i < count; i++ {
		log.Debugf("s[%s] info", i)
		log.Debugf("%s", servers[i].Info())
	}
	time.Sleep(1000 * time.Second)
	return
}
