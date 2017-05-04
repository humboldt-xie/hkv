package main

import (
	log "github.com/Sirupsen/logrus"
	kvproto "github.com/humboldt-xie/hkv/proto"

	"sync"
)

type Exporter interface {
	Start(name string, importer Importer) error
}

type ExporterManage struct {
	mu       sync.Mutex
	Exporter map[string]*DatasetExporter
	Addr     string
}

func (exp *ExporterManage) Init(config *Dataset) {
	exp.Exporter = make(map[string]*DatasetExporter)
}
func (s *ExporterManage) MakeExporter(dataset *Dataset) *DatasetExporter {
	s.mu.Lock()
	defer s.mu.Unlock()
	//mreq.Dataset["]
	name := dataset.Name
	if _, ok := s.Exporter[name]; !ok {
		if _, ok := s.Exporter[name]; !ok {
			s.Exporter[name] = &DatasetExporter{set: dataset, Addr: s.Addr}
		}
	}
	Exporter := s.Exporter[name]

	return Exporter
}

func (s *ExporterManage) GetExporter(dataset string) *DatasetExporter {
	s.mu.Lock()
	defer s.mu.Unlock()
	//mreq.Dataset["]
	if _, ok := s.Exporter[dataset]; !ok {
		return nil
	}
	Exporter := s.Exporter[dataset]
	return Exporter
}

type DatasetExporter struct {
	set          *Dataset
	IsRunning    bool
	Addr         string
	LastSequence int64
	Status       string
	mirror       kvproto.Mirror_MirrorServer
}

func (ds *DatasetExporter) Start(req *kvproto.MirrorRequest, mirror kvproto.Mirror_MirrorServer) error {
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

func (ds *DatasetExporter) CopyTo(req *kvproto.MirrorRequest, mirror kvproto.Mirror_MirrorServer) error {
	//s.set.mirror = s
	ds.mirror = mirror
	log.Debugf("start copy to %s->%s", ds.Addr, req.Addr)
	//go ds.ImportFrom(req.Addr)
	//req.Dataset.Sequence
	return ds.set.CopyTo(req.Dataset.Sequence, ds)
}

func (ds *DatasetExporter) SetStatus(status string) error {
	//ds.Status = status
	//ds.mirror.
	ds.Status = status
	return nil
}

func (ds *DatasetExporter) Copy(data *kvproto.Data) error {
	log.Debugf("exporter copy[", ds.Addr, "](", data.Sequence, ")", string(data.Key), "=>", string(data.Value))
	ds.LastSequence = data.Sequence
	return ds.mirror.Send(&kvproto.MirrorResponse{Dataset: ds.set.Name, Cmd: "copy", Data: data})
}

func (ds *DatasetExporter) Sync(data *kvproto.Data) error {
	log.Debugf("exporter sync[", ds.Addr, "](", data.Sequence, ")", string(data.Key), "=>", string(data.Value))
	ds.LastSequence = data.Sequence
	return ds.mirror.Send(&kvproto.MirrorResponse{Dataset: ds.set.Name, Cmd: "sync", Data: data})
}
