package main

import (
	"time"

	kvproto "github.com/humboldt-xie/hkv/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"fmt"
	"log"
	"sync"
)

type Importer interface {
	Dataset() *Dataset
	StartCopy() error
	EndCopy() error
	StartSync() error
	EndSync() error
	Copy(*kvproto.Data) error
	Sync(*kvproto.Data) error
	Stop() error
}

type ImporterStarter interface {
	Start(name string, importer Importer) error
}

type DatasetImporter struct {
	Name      string
	IsRunning bool
	set       *Dataset
	starter   ImporterStarter
	event     chan bool
}

func (imp *DatasetImporter) Attach() {
	select {
	case imp.event <- true:
	default:
	}
}

func (imp *DatasetImporter) Run() {
	log.Printf("importer run %s start", imp.Name)
	if imp.event != nil {
		return
	}
	imp.event = make(chan bool)
	//imp.event <- true
	go imp.Attach()
	log.Printf("importer run %s event", imp.Name)
	for {
		event := true
		if imp.IsRunning {
			event = <-imp.event
		}
		log.Printf("importer run %s", imp.Name)
		if !event {
			break
		}
		imp.IsRunning = true
		err := imp.starter.Start(imp.Name, imp)
		if err != nil {
			imp.IsRunning = false
			log.Printf("[%s]start importer error:%s", imp.Name, err)
		}

		time.Sleep(time.Second)
	}

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

/*func (imp *DatasetImporter) Do(resp *kvproto.MirrorResponse) error {
	set := imp.set
	if resp.Cmd == "copy" {
		log.Printf("copy from %s : %#v %#v", imp.client.RemoteAddr, resp.Data, set)
		if set != nil {
			set.Copy(resp.Data)
		}
	}
	if resp.Cmd == "sync" {
		log.Printf("sync from %s : %#v %#v", imp.client.RemoteAddr, resp.Data, set)
		if set != nil {
			set.Sync(resp.Data)
		}
	}
	log.Printf("[%s]cmd:%s", resp.Dataset, resp.Cmd)
	return nil
}*/
func (imp *DatasetImporter) Stop() error {
	imp.IsRunning = false
	return nil
}

type ImporterManage struct {
	mu         sync.Mutex
	RemoteAddr map[string]string
	Clients    map[string]*MirrorClient
	Importer   map[string]*DatasetImporter
	Addr       string
}

func (im *ImporterManage) Init(config *Dataset) {
	im.RemoteAddr = make(map[string]string)
	im.Clients = make(map[string]*MirrorClient)
	im.Importer = make(map[string]*DatasetImporter)
}
func (im *ImporterManage) newImporter(dataset *Dataset) *DatasetImporter {
	im.mu.Lock()
	defer im.mu.Unlock()
	//mreq.Dataset["]
	name := dataset.Name
	if _, ok := im.Importer[name]; !ok {
		im.Importer[name] = &DatasetImporter{Name: name, starter: im, set: dataset}
		go im.Importer[name].Run()
	}
	return im.Importer[name]
}

func (im *ImporterManage) getImporter(name string) *DatasetImporter {
	im.mu.Lock()
	defer im.mu.Unlock()
	//mreq.Dataset["]
	if _, ok := im.Importer[name]; !ok {
		return nil
	}
	return im.Importer[name]
}

func (im *ImporterManage) Import(dataset *Dataset, addr string) error {
	im.RemoteAddr[dataset.Name] = addr
	im.newImporter(dataset).Attach()
	return nil
}
func (im *ImporterManage) Start(name string, importer Importer) error {
	client := im.GetClient(name)
	if client == nil {
		return fmt.Errorf("client not found")
	}
	err := client.Add(name, importer)
	if err != nil {
		return err
	}
	return nil
}

func (im *ImporterManage) GetClient(name string) *MirrorClient {
	im.mu.Lock()
	defer im.mu.Unlock()
	remote := im.RemoteAddr[name]
	log.Printf("getclient[%s] %s", name, remote)
	if remote == "" {
		return nil
	}
	if _, ok := im.Clients[remote]; !ok {
		im.Clients[remote] = &MirrorClient{LocalAddr: im.Addr, RemoteAddr: remote}
		im.Clients[remote].Init()
		go im.Clients[remote].Run()
	}
	return im.Clients[remote]
}

//========================
type MirrorClient struct {
	s  *Server
	mu sync.Mutex
	//Dataset    map[string]*Dataset
	client     kvproto.Mirror_MirrorClient
	Importers  map[string]Importer
	IsRunning  bool
	LocalAddr  string
	RemoteAddr string
}

func (imp *MirrorClient) Init() {
	imp.Importers = make(map[string]Importer)
}

func (imp *MirrorClient) Add(Name string, importer Importer) error {
	imp.mu.Lock()
	defer imp.mu.Unlock()
	set := importer.Dataset()
	if imp.client == nil || !imp.IsRunning {
		return fmt.Errorf("client %s not running", imp.RemoteAddr)
	}
	err := imp.client.Send(&kvproto.MirrorRequest{Cmd: "mirror", Addr: imp.LocalAddr, Dataset: &kvproto.Dataset{Name: string(set.Name), Sequence: set.Sequence}})
	if err == nil {
		imp.Importers[Name] = importer
	}
	log.Printf("Importer add dataset:%s", Name)
	return err
}

func (imp *MirrorClient) Run() error {
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
	imp.client, err = c.Mirror(context.Background())
	if err != nil {
		return err
	}
	/*for _, v := range imp.Dataset {
		mirrorClient.Send(&kvproto.MirrorRequest{Cmd: "mirror", Addr: imp.LocalAddr, Dataset: &kvproto.Dataset{Name: string(v.Name), Sequence: v.Sequence}})
	}*/
	for {
		resp, err := imp.client.Recv()
		if err != nil {
			return err
		}
		log.Printf("mirrorclient[%s] %s %#v", imp.LocalAddr, resp.Dataset, resp.Data)
		if importer, ok := imp.Importers[resp.Dataset]; ok {
			switch resp.Cmd {
			case "copy":
				importer.Copy(resp.Data)
			case "sync":
				importer.Sync(resp.Data)
			default:
				log.Printf("mirrorclient[%s] unkown cmd %s %#v", imp.LocalAddr, resp.Dataset, resp.Cmd)

			}
			//importer.Do(resp)
		} else {
			log.Printf("mirrorclient[%s] %s %#v dataset not fount ", imp.LocalAddr, resp.Dataset, resp.Data)

		}
	}
}
