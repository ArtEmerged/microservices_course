package main

import (
	"context"
	"fmt"

	"github.com/ArtEmerged/microservices_course/week1/grpc/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address       = "localhost:5051"
	nodeID  int64 = 42
)

func main() {

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	c := client{note_v1.NewNoteV1Client(conn)}
	resp, err := c.noteServer.Get(context.Background(), &note_v1.GetRequest{NoteId: nodeID})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.GetNote())
}

type client struct {
	noteServer note_v1.NoteV1Client
}
