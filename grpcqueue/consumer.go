package grpcqueue

import (
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/ppg/grpc-queue/grpcqueue/proto"
)

// START CONSUMER DEFINITION OMIT
type service struct {
	server interface{} // the server for service methods
	md     map[string]*grpc.MethodDesc
	sd     map[string]*grpc.StreamDesc
}

// Consumer is a gRPC consumer to process gRPC requests on a queue.
type Consumer struct {
	mu       sync.Mutex
	services map[string]*service
}

// NewConsumer creates a gRPC consumer which has no services
// registered and has not started processing requests yet.
func NewConsumer() *Consumer {
	return &Consumer{
		services: make(map[string]*service),
	}
}

// END CONSUMER DEFINITION OMIT

// START CONSUME PART1 OMIT
func (c *Consumer) Consume(queue <-chan []byte) error {
	// Ingress from queue until closed
	for data := range queue {
		// Decode item
		var item pb.QueueItem
		err := proto.Unmarshal(data, &item)
		if err != nil {
			grpclog.Printf("grpcqueue: failed to unmarshal item: %s", err)
			continue
		}

		// Lookup service and method
		srv, ok := c.services[item.Service]
		if !ok {
			grpclog.Printf("grpcqueue: unknown service %v", item.Service)
			continue
		}
		md, ok := srv.md[item.Method]
		if !ok {
			grpclog.Printf("grpcqueue: unknown method %v", item.Service)
			continue
		}
		// END CONSUME PART1 OMIT

		// START CONSUME PART2 OMIT
		// Create unmarshaller
		df := func(v interface{}) error {
			if err := proto.Unmarshal(item.Payload, v.(proto.Message)); err != nil {
				return err
			}
			return nil
		}

		// Send to handler
		ctx := context.Background()
		_, err = md.Handler(srv.server, ctx, df)
		code := grpc.Code(err)
		if err != nil {
			grpclog.Printf("grpcqueue: %s - %s.%s - %s", code.String(), item.Service, item.Method, err)
			continue
		}
	}

	// Done processing return nil
	return nil
}

// END CONSUME PART2 OMIT

// RegisterService registers a service and its implementation to the gRPC
// consumer. Called from IDL generated code. This must be called before
// invoking Consume.
// START REGISTER SERVICE OMIT
func (c *Consumer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	ht := reflect.TypeOf(sd.HandlerType).Elem()
	st := reflect.TypeOf(ss)
	if !st.Implements(ht) {
		grpclog.Fatalf("grpcqueue: Consumer.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.services[sd.ServiceName]; ok {
		grpclog.Fatalf("grpcqueue: Consumer.RegisterService found duplicate service registration for %q", sd.ServiceName)
	}
	srv := &service{server: ss, md: make(map[string]*grpc.MethodDesc), sd: make(map[string]*grpc.StreamDesc)}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		srv.md[d.MethodName] = d
	}
	for i := range sd.Streams {
		d := &sd.Streams[i]
		srv.sd[d.StreamName] = d
	}
	c.services[sd.ServiceName] = srv
}

// END REGISTER SERVICE OMIT
