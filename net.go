package main

import (
	log "github.com/Sirupsen/logrus"
	kvproto "github.com/humboldt-xie/hkv/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"

	"errors"
	"fmt"
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
	Dataset  map[string]*Dataset
	Importer ImporterManage
	Exporter ExporterManage
	//private
	config  *Dataset
	db      *leveldb.DB
	cluster *ClusterServer
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
	s.cluster = &ClusterServer{current: s}
	s.cluster.Init()

	//init exporter and so
	s.Dataset = make(map[string]*Dataset)
	//s.Importer = make(map[string]*Importer)
	s.Exporter.Init(s.config)

	s.Importer.Addr = s.Addr
	s.Importer.Init(s.config)
}
func (s *Server) DbHandle() *leveldb.DB {
	return s.db
}

func (s *Server) ListenAndServe(addr string) {
	s.Addr = addr
	s.Importer.Addr = s.Addr
	s.Exporter.Addr = s.Addr
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

func (s *Server) ImportFrom(dataset string, addr string) error {
	ds := s.getDataset(dataset)
	if ds == nil {
		return ErrDatasetNotFound
	}
	return s.Importer.Import(ds, addr)
	//imp := s.getImporter(addr)
	///*for k, v := range s.Dataset {
	//	imp.Add(k, v)
	//}*/
	//imp.Add(dataset, ds)
	//go imp.ImportFrom()
	//return nil
}

// mirror server
func (s *Server) Mirror(mirror kvproto.Mirror_MirrorServer) error {
	md, _ := metadata.FromContext(mirror.Context())
	log.Debugf("mirror context :%#v", md)
	for {
		request, err := mirror.Recv()
		if err != nil {
			panic(err)
			return err
		}
		dataset := s.getDataset(request.Dataset.Name)
		if dataset == nil {
			//TODO fix error
			log.Debugf("not found exporter :%s", request.Dataset.Name)
			continue
		}

		exporter := s.Exporter.MakeExporter(dataset)
		if exporter == nil {
			//TODO fix error
			log.Debugf("not found exporter :%s", request.Dataset.Name)
			continue
		}
		go exporter.Start(request, mirror)
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
