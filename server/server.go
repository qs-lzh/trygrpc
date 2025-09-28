package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/qs-lzh/trygrpc/definition"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc"
)

type student struct {
	*pb.Student
}

type studentServiceServer struct {
	pb.UnimplementedStudentServiceServer
	students []*pb.Student
}

func (s *studentServiceServer) AddStudent(ctx context.Context, student *pb.Student) (*pb.Error, error) {
	err := s.addStudentInternal(student)
	if err != nil {
		return &pb.Error{Success: false,
			Msg: fmt.Sprint("failed to add student")}, err
	}
	return &pb.Error{}, nil
}

func (s *studentServiceServer) DescribeStudent(ctx context.Context, req *pb.StudentRequestById) (*pb.Student, error) {
	id := req.Id
	student, err := s.getStudent(int(id))
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (s *studentServiceServer) ListStudents(req *emptypb.Empty, stream grpc.ServerStreamingServer[pb.Student]) error {
	for _, student := range s.students {
		err := stream.Send(student)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *studentServiceServer) AddSeveralStudent(stream grpc.ClientStreamingServer[pb.Student, pb.Error]) error {
	for {
		student, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.Error{Success: true, Msg: "all students added"})
		}
		if err != nil {
			return err
		}
		err = s.addStudentInternal(student)
		if err != nil {
			return err
		}
	}
}

func (s *studentServiceServer) DescribeSeveralStudents(stream grpc.BidiStreamingServer[pb.StudentRequestById, pb.Student]) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		id := in.Id
		student, err := s.getStudent(int(id))
		if err != nil {
			return err
		}
		err = stream.Send(student)
		if err != nil {
			return err
		}
	}
}

func (s *studentServiceServer) addStudentInternal(student *pb.Student) error {
	s.students = append(s.students, student)
	return nil
}

func createSomeStudents() []*pb.Student {
	var students []*pb.Student
	for i := range 1000 {
		id := i + 1
		student := &pb.Student{
			Id:     int32(id),
			Name:   fmt.Sprintf("student%d", id),
			Age:    18,
			Gender: "man",
		}
		students = append(students, student)
	}
	return students
}

func (s *studentServiceServer) getStudent(id int) (*pb.Student, error) {
	for _, student := range s.students {
		if int(student.Id) == id {
			return student, nil
		}
	}
	return nil, fmt.Errorf("failed to getStudent")
}

func newStudentServiceServer() *studentServiceServer {
	return &studentServiceServer{
		students: createSomeStudents(),
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:4040")
	if err != nil {
		log.Fatal("fdsavddcvcvx exit process")
	}

	grpcServer := grpc.NewServer()
	studentServer := newStudentServiceServer()
	pb.RegisterStudentServiceServer(grpcServer, studentServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("fatal occur")
	}
}
