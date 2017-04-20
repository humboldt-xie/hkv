package main

import (
	kvproto "github.com/humboldt-xie/hkv/proto"
)

type Kv interface {
	Set(key []byte, value []byte)
	Get(key []byte) (value []byte, err error)
}

type Mirror interface {
	SetStatus(string) error
	Copy(*kvproto.Data) error
	Sync(*kvproto.Data) error
}
