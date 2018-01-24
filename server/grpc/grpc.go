package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/stphivos/todo-api-go-grpc/database"
	"github.com/stphivos/todo-api-go-grpc/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Runner ...
type Runner struct {
	Config   *models.Config
	Database database.Handler
}

// NewRunner ...
func NewRunner(config *models.Config) (*Runner, error) {
	db, err := database.Create(config)
	runner := &Runner{
		Config:   config,
		Database: db,
	}
	return runner, err
}

// Start ...
func (srv *Runner) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", srv.Config.Server.Host, srv.Config.Server.Port))
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	RegisterTodosServer(server, srv)

	err = server.Serve(listener)

	return err
}

// GetTodos ...
func (srv *Runner) GetTodos(ctx context.Context, req *Request) (*Response, error) {
	todos, err := srv.Database.GetTodos()
	if err != nil {
		// Log error but don't stop the server
		log.Println(err)

		return nil, err
	}

	res := &Response{
		Todos: srv.mapTodos(todos...),
	}

	log.Println(req.Token, ":", res)
	return res, err
}

func (srv *Runner) mapTodos(todos ...models.Todo) []*Response_Todo {
	grpcTodos := []*Response_Todo{}
	for _, todo := range todos {
		grpcTodos = append(grpcTodos, &Response_Todo{
			Id:       todo.ID.Hex(),
			Title:    todo.Title,
			Tag:      todo.Tag,
			Priority: todo.Priority,
		})
	}
	return grpcTodos
}
