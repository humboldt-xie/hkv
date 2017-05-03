package main

import (
	kvproto "github.com/humboldt-xie/hkv/proto"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)

type Cluster struct {
	current *Server
	Servers map[string]*kvproto.ServerInfo
}

func (c *Cluster) Load() error {
	//f, _ := os.OpenFile(file, os.O_RDONLY, 0)
	//TODO load from file
	d := c.current.config.ConfigGet("cluster")
	return yaml.Unmarshal([]byte(d), c)
}
func (c *Cluster) Save() error {
	//TODO Save to file
	d, err := yaml.Marshal(&c)
	if err != nil {
		panic(err)
	}
	c.current.config.ConfigSet("cluster", string(d))
	return nil
}

func (c *Cluster) Init() error {
	c.Servers = make(map[string]*kvproto.ServerInfo)
	c.Load()
	c.Update()
	c.Save()
	return nil
}

func (c *Cluster) Update() error {
	id := c.current.Id
	if c.Servers[id] == nil {
		c.Servers[id] = &kvproto.ServerInfo{Id: id, Version: 0, Addr: c.current.Addr, Dataset: make(map[string]string)}
	}
	return nil

}

func (c *Cluster) AddServer(context.Context, *kvproto.AddServerRequest) (*kvproto.AddServerReply, error) {
	return nil, nil
}

func (c *Cluster) AddDataset(ctx context.Context, req *kvproto.AddDatasetRequest) (*kvproto.AddDatasetReply, error) {
	c.current.AddDataset(req.Dataset)
	return &kvproto.AddDatasetReply{}, nil
}
func (c *Cluster) DelDataset(context.Context, *kvproto.DelDatasetRequest) (*kvproto.DelDatasetReply, error) {
	return nil, nil
}

func (c *Cluster) SyncServer(context.Context, *kvproto.SyncServerRequest) (*kvproto.SyncServerReply, error) {
	return nil, nil
}
