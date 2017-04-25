package main

import (
	kvproto "github.com/humboldt-xie/hkv/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var (
	ErrDatasetNotFound = errors.New("dataset not found")
)

type Server struct {
	Id       string
	Addr     string
	DBName   string
	Exporter map[string]*Exporter
	Importer map[string]*Importer
	Dataset  map[string]*Dataset
	//private
	config  *Dataset
	db      *leveldb.DB
	cluster *Cluster
	mu      sync.Mutex
}

func (s *Server) Init(DBName string) {

	// init db
	s.DBName = DBName
	db, err := leveldb.OpenFile(DBName, nil)
	if err != nil {
		panic(err)
	}
	s.db = db

	//init config
	s.config = &Dataset{dbHandle: s.DbHandle, Name: "config", Status: STATUS_NODE, Sequence: 0, MaxBinlog: 0}
	s.config.Init()

	s.Id = s.config.ConfigGet("serverid")
	if s.Id == "" {
		f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
		b := make([]byte, 16)
		f.Read(b)
		f.Close()
		s.Id = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		s.config.ConfigSet("serverid", s.Id)
	}

	//init cluster
	s.cluster = &Cluster{current: s, Servers: make(map[string]*kvproto.ServerInfo)}
	s.cluster.Init()

	//init exporter and so
	s.Exporter = make(map[string]*Exporter)
	s.Dataset = make(map[string]*Dataset)
	s.Importer = make(map[string]*Importer)
}
func (s *Server) DbHandle() *leveldb.DB {
	return s.db
}

func (s *Server) ListenAndServe(addr string) {
	s.Addr = addr
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	kvproto.RegisterMirrorServer(server, s)
	kvproto.RegisterClusterServer(server, s.cluster)
	server.Serve(lis)
}
func (s *Server) AddDataset(Name string) *Dataset {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ds, ok := s.Dataset[Name]; ok {
		return ds
	}
	s.Dataset[Name] = &Dataset{dbHandle: s.DbHandle, Name: Name, Status: STATUS_NODE, Sequence: 0, MaxBinlog: 10000000}
	s.Dataset[Name].Init()
	return s.Dataset[Name]
}

func (s *Server) getDataset(Name string) *Dataset {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ds, ok := s.Dataset[Name]; ok {
		return ds
	}
	return nil
	/*s.Dataset[Name] = &Dataset{dbHandle: s.DbHandle, Name: Name, Status: STATUS_NODE, Sequence: 0, MaxBinlog: 10000000}
	s.Dataset[Name].Init()
	return s.Dataset[Name]*/
}

func (s *Server) getExporter(dataset string) *Exporter {
	s.mu.Lock()
	defer s.mu.Unlock()
	//mreq.Dataset["]
	if _, ok := s.Exporter[dataset]; !ok {
		s.mu.Unlock()
		ds := s.getDataset(dataset)
		if ds == nil {
			return nil
		}
		s.mu.Lock()
		if _, ok := s.Exporter[dataset]; !ok {
			s.Exporter[dataset] = &Exporter{set: ds, Addr: s.Addr}
		}
	}
	Exporter := s.Exporter[dataset]
	return Exporter
}

func (s *Server) Set(req *kvproto.SetRequest) (*kvproto.SetReply, error) {
	ds := s.getDataset(req.Dataset)
	if ds == nil {
		return nil, ErrDatasetNotFound
	}
	return ds.Set(req)
}

func (s *Server) Get(req *kvproto.GetRequest) (*kvproto.GetReply, error) {
	ds := s.getDataset(req.Dataset)
	if ds == nil {
		return nil, ErrDatasetNotFound
	}
	return ds.Get(req)
}

func (s *Server) getImporter(addr string) *Importer {
	s.mu.Lock()
	defer s.mu.Unlock()
	//mreq.Dataset["]
	if _, ok := s.Importer[addr]; !ok {
		s.Importer[addr] = &Importer{s: s, LocalAddr: s.Addr, RemoteAddr: addr}
	}
	return s.Importer[addr]

}

func (s *Server) ImportFrom(addr string) error {
	imp := s.getImporter(addr)
	dataset := s.getDataset("1-")
	if dataset == nil {
		return ErrDatasetNotFound
	}
	imp.Add("1-", dataset)
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
		exporter := s.getExporter(request.Dataset.Name)
		if exporter == nil {
			//TODO fix error
			log.Printf("not found exporter :%s", request.Dataset.Name)
			continue
		}
		go exporter.CopyTo(request, mirror)
	}
	return nil
}

func (s *Server) Info() string {
	//res := ""
	d, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatalf("error: %v", err)
		return ""
	}
	return string(d)
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

func (imp *Importer) Add(Name string, dataset *Dataset) {
	imp.mu.Lock()
	defer imp.mu.Unlock()
	if imp.dataset == nil {
		imp.dataset = make(map[string]*Dataset)
	}
	imp.dataset[Name] = dataset
	log.Printf("Importer add dataset:%s", Name)
}

func (imp *Importer) ImportFrom() error {
	if imp.IsRunning {
		log.Printf("Importer is running %s", imp.RemoteAddr)
		return nil
	}
	imp.IsRunning = true
	defer func() {
		imp.IsRunning = false
	}()
	conn, err := grpc.Dial(imp.RemoteAddr, grpc.WithInsecure())
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
	for _, v := range imp.dataset {
		mirrorClient.Send(&kvproto.MirrorRequest{Cmd: "mirror", Addr: imp.LocalAddr, Dataset: &kvproto.Dataset{Name: string(v.Name), Sequence: v.Sequence}})
	}
	for {
		resp, err := mirrorClient.Recv()
		if err != nil {
			return err
		}
		if resp.Cmd == "copy" {
			set := imp.dataset[resp.Dataset]
			log.Printf("copy from %s : %#v %#v", imp.RemoteAddr, resp.Data, set)
			if set != nil {
				set.Copy(resp.Data)
			}
		}
		if resp.Cmd == "sync" {
			set := imp.dataset[resp.Dataset]
			log.Printf("sync from %s : %#v %#v", imp.RemoteAddr, resp.Data, set)
			if set != nil {
				set.Sync(resp.Data)
			}
		}
	}

}

type Exporter struct {
	set          *Dataset
	IsRunning    bool
	Addr         string
	LastSequence int64
	Status       string
	mirror       kvproto.Mirror_MirrorServer
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
	ds.Status = status
	return nil
}

func (ds *Exporter) Copy(data *kvproto.Data) error {
	log.Print("copy[", "](", data.Sequence, ")", string(data.Key), "=>", string(data.Value))
	ds.LastSequence = data.Sequence
	return ds.mirror.Send(&kvproto.MirrorResponse{Dataset: ds.set.Name, Cmd: "copy", Data: data})
}

func (ds *Exporter) Sync(data *kvproto.Data) error {
	log.Print("sync[", "](", data.Sequence, ")", string(data.Key), "=>", string(data.Value))
	ds.LastSequence = data.Sequence
	return ds.mirror.Send(&kvproto.MirrorResponse{Dataset: ds.set.Name, Cmd: "sync", Data: data})
}
