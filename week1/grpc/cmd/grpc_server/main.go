package main

import (
	"context"
	"fmt"
	"net"

	desc "github.com/ArtEmerged/microservices_course/week1/grpc/pkg/note_v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 5051

type server struct {
	desc.UnimplementedNoteV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{})
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}

func (s *server) Get(ctx context.Context, in *desc.GetRequest) (*desc.GetResponse, error) {
	fmt.Println("Get", "NoteID", in.GetNoteId())
	return &desc.GetResponse{
		Note: &desc.Note{
			Id: in.GetNoteId(),
			Info: &desc.NoteInfo{
				Title:    gofakeit.BeerName(),
				Content:  gofakeit.BeerMalt(),
				Author:   gofakeit.Name(),
				IsPublic: gofakeit.Bool(),
			},
			CretedAt:  timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
			DeletedAt: nil,
		},
	}, nil

}
