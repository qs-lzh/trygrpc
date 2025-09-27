package main

import (
	"context"
	"log"
	"net"

	pb "github.com/qs-lzh/trygrpc/definition"
	"google.golang.org/grpc"
)

type student struct {
	*pb.Student
}

type studentServiceServer struct {
	pb.UnimplementedStudentServiceServer
}

func (s *studentServiceServer) DescribeStudent(context.Context, *pb.StudentRequestById) (*pb.Student, error) {
	student := &pb.Student{
		Id:     2,
		Name:   "lzh",
		Age:    18,
		Gender: "man",
	}
	return student, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:4040")
	if err != nil {
		log.Fatal("fdsavddcvcvx exit process")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStudentServiceServer(grpcServer, &studentServiceServer{})

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("fdsjvhcx errrrrrrr")
	}
}
