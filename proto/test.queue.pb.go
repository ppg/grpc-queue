package proto

import (
	"fmt"

	"github.com/ppg/grpc-queue/grpcqueue"
	"golang.org/x/net/context"
)

// START REGISTER OMIT
func RegisterTestQueueConsumer(s *grpcqueue.Consumer, srv TestServer) {
	s.RegisterService(&_Test_serviceDesc, srv)
}

// END REGISTER OMIT

// START PRODUCER OMIT
// Mimics `<Service>Client` generated interface for gRPC
type TestQueueProducer interface {
	EnqueueTestRPC(ctx context.Context, in *TestRPCRequest) error
}

type testQueueProducer struct {
	queue chan<- []byte
}

func NewTestQueueProducer(queue chan<- []byte) TestQueueProducer {
	return &testQueueProducer{queue: queue}
}

func (c *testQueueProducer) EnqueueTestRPC(ctx context.Context, in *TestRPCRequest) error {
	err := grpcqueue.Enqueue(ctx, "proto.Test", "TestRPC", in, c.queue)
	if err != nil {
		return fmt.Errorf("failed to enqueue request: %s", err)
	}
	return nil
}

// END PRODUCER OMIT
