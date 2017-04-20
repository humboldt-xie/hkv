package main

import (
	kvproto "github.com/humboldt-xie/hkv/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type Server struct {
	mu       sync.Mutex
	Addr     string
	config   *Dataset
	DB       *leveldb.DB
	dss      map[string]*Exporter
	importer map[string]*Importer
	ds       map[string]*Dataset
}

func (s *Server) Init(DBName string) {
	db, err := leveldb.OpenFile(DBName, nil)
	if err != nil {
		panic(err)
	}
	s.DB = db
	s.dss = make(map[string]*Exporter)
	s.ds = make(map[string]*Dataset)
	s.importer = make(map[string]*Importer)
}
func (s *Server) DbHandle() *leveldb.DB {
	return s.DB
}

func (s *Server) ListenAndServe(addr string) {
	s.Addr = addr
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	kvproto.RegisterMirrorServer(server, s)
	server.Serve(lis)
}

func (s *Server) getDataset(Name string) *Dataset {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ds, ok := s.ds[Name]; ok {
		return ds
	}
	s.ds[Name] = &Dataset{DbHandle: s.DbHandle, Name: Name, Status: STATUS_NODE, Sequence: 0, MaxBinlog: 10000000}
	s.ds[Name].Init()
	return s.ds[Name]
}

func (s *Server) getExporter(dataset string) *Exporter {
	s.mu.Lock()
	defer s.mu.Unlock()
	//mreq.Dataset["]
	if _, ok := s.dss[dataset]; !ok {
		s.mu.Unlock()
		ds := s.getDataset(dataset)
		s.mu.Lock()
		if _, ok := s.dss[dataset]; !ok {
			s.dss[dataset] = &Exporter{set: ds, Addr: s.Addr}
		}
	}
	dss := s.dss[dataset]
	return dss
}
func (s *Server) Set(key []byte, value []byte) error {
	ds := s.getDataset("1-")
	return ds.Set(key, value)
}
func (s *Server) Get(key []byte) ([]byte, error) {
	ds := s.getDataset("1-")
	return ds.Get(key)
}

func (s *Server) getImporter(addr string) *Importer {
	s.mu.Lock()
	defer s.mu.Unlock()
	//mreq.Dataset["]
	if _, ok := s.importer[addr]; !ok {
		s.importer[addr] = &Importer{s: s, LocalAddr: s.Addr, RemoteAddr: addr}
	}
	return s.importer[addr]

}

func (s *Server) ImportFrom(addr string) error {
	imp := s.getImporter(addr)
	imp.Add("1-", s.getDataset("1-"))
	go imp.ImportFrom()
	return nil
}

// mirror server
func (s *Server) Mirror(mirror kvproto.Mirror_MirrorServer) error {
	for {
		request, err := mirror.Recv()
		if err != nil {
			panic(err)
			return err
		}
		dss := s.getExporter(request.Dataset.Name)
		go dss.CopyTo(request, mirror)
	}
	return nil
}

//========================
type Importer struct {
	s          *Server
	mu         sync.Mutex
	dataset    map[string]*Dataset
	IsRunning  bool
	LocalAddr  string
	RemoteAddr string
}

func (ds *Importer) Add(Name string, dataset *Dataset) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if ds.dataset == nil {
		ds.dataset = make(map[string]*Dataset)
	}
	ds.dataset[Name] = dataset
	log.Printf("importer add dataset:%s", Name)
}

func (ds *Importer) ImportFrom() error {
	if ds.IsRunning {
		log.Printf("importer is running %s", ds.RemoteAddr)
		return nil
	}
	ds.IsRunning = true
	defer func() {
		ds.IsRunning = false
	}()
	conn, err := grpc.Dial(ds.RemoteAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	c := kvproto.NewMirrorClient(conn)
	mirrorClient, err := c.Mirror(context.Background())
	if err != nil {
		return err
	}
	for _, v := range ds.dataset {
		mirrorClient.Send(&kvproto.MirrorRequest{Cmd: "mirror", Addr: ds.LocalAddr, Dataset: &kvproto.Dataset{Name: string(v.Name), Sequence: v.Sequence}})
	}
	for {
		resp, err := mirrorClient.Recv()
		if err != nil {
			return err
		}
		if resp.Cmd == "copy" {
			set := ds.dataset[resp.Dataset]
			log.Printf("copy from %s : %#v %#v", ds.RemoteAddr, resp.Data, set)
			if set != nil {
				set.Copy(resp.Data)
			}
		}
		if resp.Cmd == "sync" {
			set := ds.dataset[resp.Dataset]
			log.Printf("sync from %s : %#v %#v", ds.RemoteAddr, resp.Data, set)
			if set != nil {
				set.Sync(resp.Data)
			}
		}
	}

}

type Exporter struct {
	set       *Dataset
	IsRunning bool
	Addr      string
	mirror    kvproto.Mirror_MirrorServer
}

func (ds *Exporter) CopyTo(req *kvproto.MirrorRequest, mirror kvproto.Mirror_MirrorServer) error {
	if ds.IsRunning {
		return nil
	}
	ds.IsRunning = true
	defer func() {
		ds.IsRunning = false
	}()
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
	return nil
}

func (ds *Exporter) Copy(data *kvproto.Data) error {
	log.Print("copy[", "](", data.Sequence, ")", string(data.Key), "=>", string(data.Value))
	return ds.mirror.Send(&kvproto.MirrorResponse{Dataset: ds.set.Name, Cmd: "copy", Data: data})
}

func (ds *Exporter) Sync(data *kvproto.Data) error {
	log.Print("sync[", "](", data.Sequence, ")", string(data.Key), "=>", string(data.Value))
	return ds.mirror.Send(&kvproto.MirrorResponse{Dataset: ds.set.Name, Cmd: "sync", Data: data})
}
