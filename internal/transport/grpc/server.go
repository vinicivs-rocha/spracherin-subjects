package grpc

import (
	"net"

	subjectsv1 "github.com/joserocha/spracherin/spracherin-proto/gen/go/github.com/vinicivs-rocha/spracherin/spracherin-proto/gen/go/subjects/v1"
	"github.com/vinicivs-rocha/spracherin-subjects/internal/data"
	"google.golang.org/grpc"
)

func NewSubjectGRPCServer(repository data.SubjectRepository, messager data.Messager) *grpc.Server {
	server := grpc.NewServer()
	handler := NewSubjectHandler(repository, messager)
	subjectsv1.RegisterSubjectServiceServer(server, handler)
	return server
}

func ListenAndServe(addr string, repository data.SubjectRepository, messager data.Messager) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := NewSubjectGRPCServer(repository, messager)
	return server.Serve(listener)
}
