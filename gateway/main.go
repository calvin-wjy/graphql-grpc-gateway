// gateway/main.go
package main

import (
	"log"

	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	pb "github.com/crshao/grpc-graphql-gateway/student"
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
)

func main() {
	mux := runtime.NewServeMux()

	if err := pb.RegisterStudentManagementGraphql(mux); err != nil {
		log.Fatalln(err)
	}

	http.Handle("/", playground.Handler("GraphQL", "/graphql"))
	http.Handle("/graphql", mux)
	log.Fatalln(http.ListenAndServe(":8888", nil))
}
