package main

import (
	"context"
	"log"
	"time"

	pb "github.com/qs-lzh/trygrpc/definition"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStudentServiceClient(conn)

	req := &pb.StudentRequestById{
		Id: 2,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.DescribeStudent(ctx, req)
	if err != nil {
		log.Fatalf("could not describe student: %v", err)
	}

	log.Printf("Student: ID=%d, Name=%s, Age=%d, Gender=%s",
		resp.Id, resp.Name, resp.Age, resp.Gender)
}
