package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/qs-lzh/trygrpc/definition"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc"
)

func main() {
	// 连接服务端
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStudentServiceClient(conn)

	// ---- 1. 调用 AddStudent ----
	student := &pb.Student{
		Id:     10001,
		Name:   "Alice",
		Age:    20,
		Gender: "woman",
	}
	_, err = client.AddStudent(context.Background(), student)
	if err != nil {
		log.Fatalf("AddStudent failed: %v", err)
	}
	fmt.Println("AddStudent success")

	// ---- 2. 调用 DescribeStudent ----
	req := &pb.StudentRequestById{Id: 1}
	stu, err := client.DescribeStudent(context.Background(), req)
	if err != nil {
		log.Fatalf("DescribeStudent failed: %v", err)
	}
	fmt.Printf("DescribeStudent success: got student %v\n", stu.Name)

	// ---- 3. 调用 ListStudents ----
	stream1, err := client.ListStudents(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("ListStudents failed: %v", err)
	}
	count := 0
	for {
		_, err := stream1.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ListStudents recv failed: %v", err)
		}
		count++
		if count >= 3 { // 只打印前几个
			break
		}
	}
	fmt.Println("ListStudents success (received at least 3 students)")

	// ---- 4. 调用 AddSeveralStudent (client streaming) ----
	stream2, err := client.AddSeveralStudent(context.Background())
	if err != nil {
		log.Fatalf("AddSeveralStudent failed: %v", err)
	}
	for i := 0; i < 3; i++ {
		stu := &pb.Student{
			Id:     int32(2000 + i),
			Name:   fmt.Sprintf("BatchStu%d", i),
			Age:    21,
			Gender: "man",
		}
		if err := stream2.Send(stu); err != nil {
			log.Fatalf("AddSeveralStudent send failed: %v", err)
		}
	}
	res, err := stream2.CloseAndRecv()
	if err != nil {
		log.Fatalf("AddSeveralStudent close failed: %v", err)
	}
	fmt.Printf("AddSeveralStudent success: %v\n", res.Msg)

	// ---- 5. 调用 DescribeSeveralStudents (bi-directional streaming) ----
	stream3, err := client.DescribeSeveralStudents(context.Background())
	if err != nil {
		log.Fatalf("DescribeSeveralStudents failed: %v", err)
	}

	// 启动一个 goroutine 发送请求
	go func() {
		for i := 1; i <= 3; i++ {
			req := &pb.StudentRequestById{Id: int32(i)}
			if err := stream3.Send(req); err != nil {
				log.Fatalf("DescribeSeveralStudents send failed: %v", err)
			}
			time.Sleep(200 * time.Millisecond)
		}
		stream3.CloseSend()
	}()

	// 接收返回
	for {
		stu, err := stream3.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("DescribeSeveralStudents recv failed: %v", err)
		}
		fmt.Printf("DescribeSeveralStudents success: got %s\n", stu.Name)
	}
}
