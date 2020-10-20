package main

import (
	"context"
	"log"

	p "github.com/Amirkin/product-service/proto"
	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.Flags() | log.Llongfile)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:5300", opts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := p.NewProductServiceClient(conn)
	resp, err := client.List(context.Background(), &p.ListParams{
		PagingParams: &p.PagingParams{
			Page:   0,
			Offset: 0,
			Limit:  0,
		},
		SortParams: &p.SortParams{
			Name:       0,
			Price:      0,
			LastUpdate: 0,
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%v \n", resp)
}
