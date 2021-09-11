// service/main.go
package main

import (
	"context"
	"log"
	"net"

	pb "github.com/crshao/grpc-graphql-gateway/student"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedStudentManagementServer
}

func (s *Server) GetStudents(ctx context.Context, req *pb.GetStudentsParams) (*pb.StudentsList, error) {
	var students_list *pb.StudentsList = &pb.StudentsList{}

	student := pb.Student{Name: "test", Nim: "123123"}
	students_list.Students = append(students_list.Students, &student)

	return students_list, nil
}

func main() {
	conn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server := grpc.NewServer()

	pb.RegisterStudentManagementServer(server, &Server{})
	log.Printf("Server listening at %v", conn.Addr())
	server.Serve(conn)
}
