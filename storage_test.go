package main

import (
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	s1 := Server{}
	s1.Init("s1")
	go s1.ListenAndServe("127.0.0.1:7001")
	time.Sleep(1 * time.Second)
	s2 := Server{}
	s2.Init("s2")
	go s2.ListenAndServe("127.0.0.1:7002")
	time.Sleep(1 * time.Second)
	//go s2.MirrorTo("127.0.0.1:7001")

	s1.Set([]byte("hello"), []byte("world"))
	time.Sleep(1 * time.Second)
	t.Fatal(s2.Get([]byte("hello")))

	//time.Sleep(2 * time.Second)
}
