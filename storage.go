package main

import (
	pb "github.com/humboldt-xie/hkv/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

var SlotCount = 16

type Kv interface {
	Set(key []byte, value []byte)
	Get(key []byte) (value []byte, err error)
}

type Mirror interface {
	SetStatus(string) error
	Copy(Item) error
	Sync(Item) error
}

type Item struct {
	Sequence int64
	Key      []byte
	Value    []byte
}

func GetDb() *leveldb.DB {
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		panic(err)
	}
	return db
}

type Node struct {
	Id      string
	Dataset map[string]*Dataset
}

func NewNode(Id string) *Node {
	return &Node{Id: Id}
}

type DatasetMirrorServer struct {
	set *Dataset
}

func (s *DatasetMirrorServer) Mirror(req *pb.MirrorRequest, mirror pb.Data_MirrorServer) error {
	//s.set.mirror = s
	go s.set.CopyTo(s)
	return nil
}

func (ds *DatasetMirrorServer) SetStatus(status string) error {
	//ds.Status = status
	return nil
}

func (ds *DatasetMirrorServer) Copy(item Item) error {
	return nil //DB.Put(ds.Key(item.Key), item.Value, nil)
}

func (ds *DatasetMirrorServer) Sync(item Item) error {
	return nil //DB.Put(ds.Key(item.Key), item.Value, nil)
}

type Service struct {
	mu  sync.Mutex
	dss map[string]*DatasetMirrorServer
	ds  map[string]*Dataset
}

func (s *Service) Init() {
	s.dss = make(map[string]*DatasetMirrorServer)
	s.ds = make(map[string]*Dataset)
	//s.GetDataset("1-")
}

func (s *Service) GetDataset(Name string) *Dataset {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ds, ok := s.ds[Name]; ok {
		return ds
	}
	s.ds[Name] = &Dataset{Name: []byte(Name), Status: STATUS_NODE, Sequence: 0}
	return s.ds[Name]
}

func (s *Service) StartMirror(mreq *pb.MirrorRequest, mirror pb.Data_MirrorServer) {
	s.mu.Lock()
	//mreq.Dataset["]
	dataset := "1-"
	if _, ok := s.dss[dataset]; !ok {
		s.dss[dataset] = &DatasetMirrorServer{set: s.GetDataset(dataset)}
	}
	dss := s.dss[dataset]
	s.mu.Unlock()
	dss.Mirror(mreq, mirror)
}

func (s *Service) Mirror(mirror pb.Data_MirrorServer) {
	for {
		request, err := mirror.Recv()
		if err != nil {
			break
		}
		s.StartMirror(request, mirror)
	}

}

/*func (nd *DataNode) Mirror(req pb.MirrorRequest, reply pb.Node_MirrorServer) error {
	shard := nd.Shards[req.ShardId]
	if shard == nil {
		return fmt.Errorf("node %s shard %d not existing", nd.NodeId, req.ShardId)
	}
	return shard.Mirror(req, reply)
}*/

var DB = GetDb()

func main() {
	DB.Close()
}
