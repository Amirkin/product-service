package main

import (
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc/grpclog"

	grpc "google.golang.org/grpc"

	p "github.com/Amirkin/product-service/proto"
)

func main() {
	log.SetFlags(log.Flags() | log.Llongfile)

	store, err := NewStore("mongodb://deploy_mongo_1:27017")
	if err != nil {
		log.Fatalln(err)
	}

	listener, err := net.Listen("tcp", ":5300")

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	handlers := &Handlers{
		Store: store,
		Requester: &Requester{
			http.DefaultClient,
		},
	}

	p.RegisterProductServiceServer(server, handlers)
	log.Fatalln(server.Serve(listener))
}
