// service/main.go
package main

import (
	"context"
	"log"
	"net"

	pb "github.com/crshao/grpc-graphql-gateway/student"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

func NewStudentManagementServer() *Server {
	return &Server{}
}

type Server struct {
	conn *pgx.Conn
	pb.UnimplementedStudentManagementServer
}

func (server *Server) Run() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStudentManagementServer(s, server)
	log.Printf("Server listening at %v", lis.Addr())

	return s.Serve(lis)
}

func (s *Server) GetStudents(ctx context.Context, req *pb.GetStudentsParams) (*pb.StudentsList, error) {
	var students_list *pb.StudentsList = &pb.StudentsList{}

	rows, err := s.conn.Query(context.Background(), "SELECT * FROM students")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		student := pb.Student{}
		err = rows.Scan(&student.Id, &student.Name, &student.Nim)
		if err != nil {
			return nil, err
		}
		students_list.Students = append(students_list.Students, &student)
	}

	return students_list, nil
}

func main() {

	database_url := "postgres://calvin.wijaya:calvin.wijaya@localhost:5432/grpc_graphql_gateway"
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}

	defer conn.Close(context.Background())

	var student_management_server *Server = NewStudentManagementServer()
	student_management_server.conn = conn

	if err := student_management_server.Run(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
