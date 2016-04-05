package main

import (
	"log"
	"time"

	"google.golang.org/grpc/grpclog"

	"golang.org/x/net/context"

	"github.com/ppg/grpc-queue/grpcqueue"
	pb "github.com/ppg/grpc-queue/proto"
)

func init() {
	// Prettier CLI output
	log.SetFlags(0)
}

// START MAIN CONSUMER OMIT
func main() {
	// Make an in memory queue
	queue := make(chan []byte, 10)

	// Create a consumer
	consumer := grpcqueue.NewConsumer()

	// Create a service implementation and register according to proto IDL
	testService := &testServer{}
	pb.RegisterTestQueueConsumer(consumer, testService)

	// Start consumer and wait for channel to close (in background)
	go func() {
		if err := consumer.Consume(queue); err != nil {
			log.Fatalf("consume failed: %s", err)
		}
		log.Fatal("consume stopped")
	}()
	// END MAIN CONSUMER OMIT

	// START MAIN PRODUCER OMIT
	// Create a producer
	producer := pb.NewTestQueueProducer(queue)

	// Enqueue a couple objects
	ctx := context.Background()
	log.Print("Enqueue: Hello World")
	producer.EnqueueTestRPC(ctx, &pb.TestRPCRequest{Message: "Hello World"})
	log.Print("Enqueue: Where am I?")
	producer.EnqueueTestRPC(ctx, &pb.TestRPCRequest{Message: "Where am I?"})
	log.Print("Enqueue: Unknown on foo.Bar")
	grpcqueue.Enqueue(ctx, "foo.Bar", "Unknown", &pb.TestRPCRequest{}, queue)
	// Wait a little for these messages
	log.Print("Waiting")
	time.Sleep(1 * time.Second)

	// Enqueue goodbye
	log.Print("Enqueue: Goodbye!")
	producer.EnqueueTestRPC(ctx, &pb.TestRPCRequest{Message: "Goodbye!"})
	// END MAIN PRODUCER OMIT

	// START MAIN WAIT OMIT
	// Close channel to exit
	close(queue)
	// Wait a little for these messages
	log.Print("Waiting")
	time.Sleep(1 * time.Second)
}

// END MAIN WAIT OMIT

// START TEST SERVER OMIT
type testServer struct {
}

func (testServer) TestRPC(ctx context.Context, req *pb.TestRPCRequest) (*pb.TestRPCResponse, error) {
	grpclog.Printf("[testServer] %s", req.Message)
	return &pb.TestRPCResponse{}, nil
}

// END TEST SERVER OMIT
