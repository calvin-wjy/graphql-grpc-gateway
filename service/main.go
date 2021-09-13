// service/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	rd "github.com/crshao/grpc-graphql-gateway/cache"
	pb "github.com/crshao/grpc-graphql-gateway/student"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

// type Student struct {

// }

var redis *rd.Client

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

	val, err := redis.Get(ctx, "sudah")
	if err != nil {
		fmt.Println("redis.get()")
		panic(err)
	}
	fmt.Println("val: ", val, "|")

	// Sudah di set sebelumnya
	if val.Val() != "" {
		// fmt.Println("\n\nSebelum di unmarshal:\n", val.Val())
		val_byte := []byte(val.Val())
		json.Unmarshal(val_byte, &students_list)
		// fmt.Println("Setelah unmarshal:\n", students_list)

		fmt.Println("AMBIL DARI REDIS")

		return students_list, nil
	}

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
	students_list_str, _ := json.Marshal(students_list)
	fmt.Printf("%s\n", string(students_list_str))

	// set
	redis.Set(ctx, "sudah", students_list)

	return students_list, nil
}

func main() {
	database_url := "postgres://calvin.wijaya:calvin.wijaya@localhost:5432/grpc_graphql_gateway"
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())

	// Redis
	redis, err = rd.NewRedis()
	if err != nil {
		log.Fatalf("Could not initialize Redis client %s", err)
	}

	var student_management_server *Server = NewStudentManagementServer()
	student_management_server.conn = conn

	if err := student_management_server.Run(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
