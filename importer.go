package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	kvproto "github.com/humboldt-xie/hkv/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"fmt"
	"sync"
)

//Importer 是数据导入接口与Exporter 一一对应
//
type Importer interface {
	Name() string
	Sequence() int64
	StartCopy() error
	EndCopy() error
	StartSync() error
	EndSync() error
	Copy(*kvproto.Data) error
	Sync(*kvproto.Data) error
	Stop() error
}

type DatasetImporter struct {
	IsRunning  bool
	set        *Dataset
	RemoteAddr string
	LocalAddr  string
}

func (imp *DatasetImporter) start() {
	imp.IsRunning = true
}
func (imp *DatasetImporter) stop() {
	imp.IsRunning = false
}

func (imp *DatasetImporter) Name() string {
	return imp.set.Name
}
func (imp *DatasetImporter) Sequence() int64 {
	return imp.set.Sequence
}
func (imp *DatasetImporter) Dataset() *Dataset {
	return imp.set
}
func (imp *DatasetImporter) StartCopy() error {
	//TODO start copy
	return nil
}
func (imp *DatasetImporter) EndCopy() error {
	//TODO end copy
	return nil
}
func (imp *DatasetImporter) StartSync() error {
	//TODO start sync
	return nil
}
func (imp *DatasetImporter) EndSync() error {
	//TODO end sync
	return nil
}

func (imp *DatasetImporter) Copy(data *kvproto.Data) error {
	return imp.set.Copy(data)
}

func (imp *DatasetImporter) Sync(data *kvproto.Data) error {
	return imp.set.Sync(data)
}
func (imp *DatasetImporter) Run() {
	if imp.IsRunning {
		return
	}
	imp.start()
	defer imp.stop()
	for {
		imp.run()
		time.Sleep(1 * time.Second)
	}
}

func (imp *DatasetImporter) run() {
	rpcClient, err := rpcPool.GetClient(imp.RemoteAddr)
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return
	}
	c := kvproto.NewMirrorClient(rpcClient.Conn)
	client, err := c.Mirror(context.Background())
	if err != nil {
		log.Errorf("mirrorclient[%s] %s %s create client error", imp.LocalAddr, imp.RemoteAddr, err)
		return
	}
	defer client.CloseSend()
	err = client.Send(&kvproto.MirrorRequest{
		Cmd:  "mirror",
		Addr: imp.LocalAddr,
		Dataset: &kvproto.Dataset{
			Name:     imp.Name(),
			Sequence: imp.Sequence(),
		},
	})
	if err != nil {
		log.Errorf("mirrorclient[%s] %s %s send error", imp.LocalAddr, imp.RemoteAddr, err)
		return
	}
	for {
		resp, err := client.Recv()
		if err != nil {
			log.Errorf("mirrorclient[%s] %s %s recv error", imp.LocalAddr, imp.RemoteAddr, err)
			return
		}
		if imp.set.Name != resp.Dataset {
			log.Debugf("mirrorclient[%s] %s %#v dataset not fount ", imp.LocalAddr, resp.Dataset, resp.Data)
			return
		}
		log.Debugf("mirrorclient[%s]:%s %s %#v", imp.LocalAddr, resp.Cmd, resp.Dataset, resp.Data)
		switch resp.Cmd {
		case "copy":
			imp.Copy(resp.Data)
		case "sync":
			imp.Sync(resp.Data)
		default:
			log.Errorf("mirrorclient[%s] unkown cmd %s %#v", imp.LocalAddr, resp.Dataset, resp.Cmd)
		}
	}

}

//func (imp *DatasetImporter) Do(resp *kvproto.MirrorResponse) error {
//	set := imp.set
//	if resp.Cmd == "copy" {
//		log.Debugf("copy from %s : %#v %#v", imp.client.RemoteAddr, resp.Data, set)
//		if set != nil {
//			set.Copy(resp.Data)
//		}
//	}
//	if resp.Cmd == "sync" {
//		log.Debugf("sync from %s : %#v %#v", imp.client.RemoteAddr, resp.Data, set)
//		if set != nil {
//			set.Sync(resp.Data)
//		}
//	}
//	log.Debugf("[%s]cmd:%s", resp.Dataset, resp.Cmd)
//	return nil
//}
func (imp *DatasetImporter) Stop() error {
	imp.IsRunning = false
	return nil
}

type ImporterManage struct {
	mu sync.Mutex
	//RemoteAddr map[string]string
	//Clients    map[string]*MirrorClient
	Importer map[string]*DatasetImporter
	Addr     string
}

func (im *ImporterManage) Init(config *Dataset) {
	//im.RemoteAddr = make(map[string]string)
	//im.Clients = make(map[string]*MirrorClient)
	im.Importer = make(map[string]*DatasetImporter)
}

func (im *ImporterManage) Import(dataset *Dataset, addr string) error {
	//im.RemoteAddr[dataset.Name] = addr
	//im.newImporter(dataset).Attach()
	im.mu.Lock()
	defer im.mu.Unlock()
	//mreq.Dataset["]
	name := addr + "->" + dataset.Name
	ds, ok := im.Importer[name]
	if !ok {
		ds = &DatasetImporter{RemoteAddr: addr, set: dataset, LocalAddr: im.Addr}
		im.Importer[name] = ds
	}
	ds.LocalAddr = im.Addr
	go ds.Run()
	return nil
}

type RpcClient struct {
	mu      sync.Mutex
	Address string
	Conn    *grpc.ClientConn
	Error   error
}

func (client *RpcClient) Reconnect() {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.Error != nil {
		client.Conn, client.Error = grpc.Dial(client.Address, grpc.WithInsecure())
	}
}

type RpcPool struct {
	mu      sync.Mutex
	clients map[string]*RpcClient
}

func (rpc *RpcPool) GetClient(address string) (rpcClient *RpcClient, err error) {
	rpc.mu.Lock()
	var ok bool
	rpcClient, ok = rpc.clients[address]
	if !ok {
		rpcClient = &RpcClient{Address: address, Conn: nil, Error: fmt.Errorf("no conn %s", address)}
		rpc.clients[address] = rpcClient
	}
	rpc.mu.Unlock()
	if rpcClient.Error != nil {
		rpcClient.Reconnect()
	}
	return rpcClient, rpcClient.Error
}

func NewRpcPool() *RpcPool {
	return &RpcPool{clients: make(map[string]*RpcClient)}
}

var rpcPool = NewRpcPool()
