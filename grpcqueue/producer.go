package grpcqueue

import (
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	pb "github.com/ppg/grpc-queue/grpcqueue/proto"
)

// Enqueue is called by the generated code. It sends the RPC request to the
// queueing service.
func Enqueue(ctx context.Context, service, method string, args interface{}, queue chan<- []byte,
) (err error) {
	// Generate QueueItem
	item := &pb.QueueItem{Service: service, Method: method}

	// Marshal message
	item.Payload, err = proto.Marshal(args.(proto.Message))
	if err != nil {
		return err
	}

	// Marshal queue item
	data, err := proto.Marshal(item)
	if err != nil {
		return err
	}

	queue <- data
	return nil
}
