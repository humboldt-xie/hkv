package main

import (
	kvproto "github.com/humboldt-xie/hkv/proto"
	"golang.org/x/net/context"
	//"gopkg.in/yaml.v2"
)

type Cluster struct {
	current *Server
	Servers map[string]*kvproto.ServerInfo
}

func (c *Cluster) Init() error {

	return nil

}

func (c *Cluster) AddServer(context.Context, *kvproto.AddServerRequest) (*kvproto.AddServerReply, error) {
	return nil, nil
}

func (c *Cluster) AddDataset(context.Context, *kvproto.AddDatasetRequest) (*kvproto.AddDatasetReply, error) {

	return nil, nil
}
func (c *Cluster) DelDataset(context.Context, *kvproto.DelDatasetRequest) (*kvproto.DelDatasetReply, error) {
	return nil, nil
}

func (c *Cluster) SyncServer(context.Context, *kvproto.SyncServerRequest) (*kvproto.SyncServerReply, error) {
	return nil, nil
}
