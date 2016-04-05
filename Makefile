proto_files = $(wildcard proto/*.proto)
pb_files = $(proto_files:.proto=.pb.go)

%.pb.go : %.proto
	protoc --go_out=plugins=grpc:./proto --proto_path=./proto $(proto_files)

grpcqueue/proto/queue_item.pb.go : grpcqueue/proto/queue_item.proto
	protoc --go_out=plugins=grpc:./grpcqueue/proto --proto_path=./grpcqueue/proto grpcqueue/proto/queue_item.proto

all: $(pb_files) grpcqueue/proto/queue_item.pb.go

clean:
	rm -f $(pb_files)
