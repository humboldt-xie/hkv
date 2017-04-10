all: proto/kv.pb.go hkv

test: *.go
	go test ./

hkv:*.go 
	go build ./

proto/kv.pb.go: proto/kv.proto
	protoc -I proto/ proto/kv.proto --go_out=plugins=grpc:proto
