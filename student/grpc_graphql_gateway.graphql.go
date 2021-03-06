// Code generated by proroc-gen-graphql, DO NOT EDIT.
package grpc_graphql_gateway

import (
	"context"

	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	gql__type_StudentsList  *graphql.Object      // message StudentsList in student.proto
	gql__type_Student       *graphql.Object      // message Student in student.proto
	gql__input_StudentsList *graphql.InputObject // message StudentsList in student.proto
	gql__input_Student      *graphql.InputObject // message Student in student.proto
)

func Gql__type_StudentsList() *graphql.Object {
	if gql__type_StudentsList == nil {
		gql__type_StudentsList = graphql.NewObject(graphql.ObjectConfig{
			Name: "GrpcGraphqlGateway_Type_StudentsList",
			Fields: graphql.Fields{
				"students": &graphql.Field{
					Type: graphql.NewList(Gql__type_Student()),
				},
			},
		})
	}
	return gql__type_StudentsList
}

func Gql__type_Student() *graphql.Object {
	if gql__type_Student == nil {
		gql__type_Student = graphql.NewObject(graphql.ObjectConfig{
			Name: "GrpcGraphqlGateway_Type_Student",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"nim": &graphql.Field{
					Type: graphql.String,
				},
				"id": &graphql.Field{
					Type: graphql.Int,
				},
			},
		})
	}
	return gql__type_Student
}

func Gql__input_StudentsList() *graphql.InputObject {
	if gql__input_StudentsList == nil {
		gql__input_StudentsList = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "GrpcGraphqlGateway_Input_StudentsList",
			Fields: graphql.InputObjectConfigFieldMap{
				"students": &graphql.InputObjectFieldConfig{
					Type: graphql.NewList(Gql__input_Student()),
				},
			},
		})
	}
	return gql__input_StudentsList
}

func Gql__input_Student() *graphql.InputObject {
	if gql__input_Student == nil {
		gql__input_Student = graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "GrpcGraphqlGateway_Input_Student",
			Fields: graphql.InputObjectConfigFieldMap{
				"name": &graphql.InputObjectFieldConfig{
					Type: graphql.String,
				},
				"nim": &graphql.InputObjectFieldConfig{
					Type: graphql.String,
				},
				"id": &graphql.InputObjectFieldConfig{
					Type: graphql.Int,
				},
			},
		})
	}
	return gql__input_Student
}

// graphql__resolver_StudentManagement is a struct for making query, mutation and resolve fields.
// This struct must be implemented runtime.SchemaBuilder interface.
type graphql__resolver_StudentManagement struct {

	// Automatic connection host
	host string

	// grpc dial options
	dialOptions []grpc.DialOption

	// grpc client connection.
	// this connection may be provided by user
	conn *grpc.ClientConn
}

// new_graphql_resolver_StudentManagement creates pointer of service struct
func new_graphql_resolver_StudentManagement(conn *grpc.ClientConn) *graphql__resolver_StudentManagement {
	return &graphql__resolver_StudentManagement{
		conn: conn,
		host: "localhost:50051",
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
	}
}

// CreateConnection() returns grpc connection which user specified or newly connected and closing function
func (x *graphql__resolver_StudentManagement) CreateConnection(ctx context.Context) (*grpc.ClientConn, func(), error) {
	// If x.conn is not nil, user injected their own connection
	if x.conn != nil {
		return x.conn, func() {}, nil
	}

	// Otherwise, this handler opens connection with specified host
	conn, err := grpc.DialContext(ctx, x.host, x.dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	return conn, func() { conn.Close() }, nil
}

// GetQueries returns acceptable graphql.Fields for Query.
func (x *graphql__resolver_StudentManagement) GetQueries(conn *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
		"getStudents": &graphql.Field{
			Type: Gql__type_StudentsList(),
			Args: graphql.FieldConfigArgument{},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var req GetStudentsParams
				if err := runtime.MarshalRequest(p.Args, &req, false); err != nil {
					return nil, errors.Wrap(err, "Failed to marshal request for getStudents")
				}
				client := NewStudentManagementClient(conn)
				resp, err := client.GetStudents(p.Context, &req)
				if err != nil {
					return nil, errors.Wrap(err, "Failed to call RPC GetStudents")
				}
				return resp, nil
			},
		},
	}
}

// GetMutations returns acceptable graphql.Fields for Mutation.
func (x *graphql__resolver_StudentManagement) GetMutations(conn *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{}
}

// Register package divided graphql handler "without" *grpc.ClientConn,
// therefore gRPC connection will be opened and closed automatically.
// Occasionally you may worry about open/close performance for each handling graphql request,
// then you can call RegisterStudentManagementGraphqlHandler with *grpc.ClientConn manually.
func RegisterStudentManagementGraphql(mux *runtime.ServeMux) error {
	return RegisterStudentManagementGraphqlHandler(mux, nil)
}

// Register package divided graphql handler "with" *grpc.ClientConn.
// this function accepts your defined grpc connection, so that we reuse that and never close connection inside.
// You need to close it maunally when application will terminate.
// Otherwise, you can specify automatic opening connection with ServiceOption directive:
//
// service StudentManagement {
//    option (graphql.service) = {
//        host: "host:port"
//        insecure: true or false
//    };
//
//    ...with RPC definitions
// }
func RegisterStudentManagementGraphqlHandler(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return mux.AddHandler(new_graphql_resolver_StudentManagement(conn))
}
