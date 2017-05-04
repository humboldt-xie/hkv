package main

import (
	"log"

	kvproto "github.com/humboldt-xie/hkv/proto"
)

type Exporter struct {
	set          *Dataset
	IsRunning    bool
	Addr         string
	LastSequence int64
	Status       string
	mirror       kvproto.Mirror_MirrorServer
}

func (ds *Exporter) Start(req *kvproto.MirrorRequest, mirror kvproto.Mirror_MirrorServer) error {
	if ds.IsRunning {
		return nil
	}
	ds.IsRunning = true
	defer func() {
		ds.IsRunning = false
	}()
	ds.mirror = mirror
	ds.mirror.Send(&kvproto.MirrorResponse{Dataset: ds.set.Name, Cmd: "start_copy", Data: nil})
	return ds.CopyTo(req, mirror)

}

func (ds *Exporter) CopyTo(req *kvproto.MirrorRequest, mirror kvproto.Mirror_MirrorServer) error {
	//s.set.mirror = s
	ds.mirror = mirror
	log.Printf("start copy to %s->%s", ds.Addr, req.Addr)
	//go ds.ImportFrom(req.Addr)
	//req.Dataset.Sequence
	return ds.set.CopyTo(req.Dataset.Sequence, ds)
}

func (ds *Exporter) SetStatus(status string) error {
	//ds.Status = status
	//ds.mirror.
	ds.Status = status
	return nil
}

func (ds *Exporter) Copy(data *kvproto.Data) error {
	log.Print("exporter copy[", ds.Addr, "](", data.Sequence, ")", string(data.Key), "=>", string(data.Value))
	ds.LastSequence = data.Sequence
	return ds.mirror.Send(&kvproto.MirrorResponse{Dataset: ds.set.Name, Cmd: "copy", Data: data})
}

func (ds *Exporter) Sync(data *kvproto.Data) error {
	log.Print("exporter sync[", ds.Addr, "](", data.Sequence, ")", string(data.Key), "=>", string(data.Value))
	ds.LastSequence = data.Sequence
	return ds.mirror.Send(&kvproto.MirrorResponse{Dataset: ds.set.Name, Cmd: "sync", Data: data})
}
