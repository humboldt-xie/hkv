all: proto/kv.pb.go  hkvd

test: *.go
	go test ./

hkvd:*.go 
	go build -o hkvd ./

proto/kv.pb.go: proto/kv.proto
	protoc -I proto/ proto/kv.proto --go_out=plugins=grpc:proto
