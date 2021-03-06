package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// go run example/grpc_client/app.go
func main() {
	// Set up a connection to the server.
	cr, err := credentials.NewClientTLSFromFile("./testdata/out/GearTest.crt", "127.0.0.1")
	if err != nil {
		log.Fatalf("credentials error: %v", err)
	}
	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithTransportCredentials(cr))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	for i := 0; i < 10; i++ {
		r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: strconv.Itoa(i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.Message)
	}

	// bench(c)
}

func bench(c pb.GreeterClient) {
	var total = 1000000
	var cocurrency = 1000
	var w sync.WaitGroup

	w.Add(total)
	co := make(chan int, cocurrency)

	task := func(i int) {
		defer w.Done()
		_, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "Ping"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		<-co
		if (i % 1000) == 0 {
			fmt.Print(".")
		}
	}

	t := time.Now()
	for i := 0; i < total; i++ {
		co <- i
		go task(i)
	}
	log.Println("Wait")
	w.Wait()
	sec := int(time.Now().Sub(t) / 1e6)
	log.Printf("\nFinished, cocurrency: %d, time: %d ms, %d ops", cocurrency, sec, (total / (sec / 1000)))
}
