package server

import (
	"fmt"

	"github.com/stphivos/todo-api-go-grpc/models"
	"github.com/stphivos/todo-api-go-grpc/server/grpc"
)

// Runner ...
type Runner interface {
	Start() error
}

// Create ...
func Create(config *models.Config) (Runner, error) {
	var srv Runner
	var err error

	switch config.Server.Type {
	case "grpc":
		srv, err = grpc.NewRunner(config)
	default:
		err = fmt.Errorf("Server type %v is not supported", config.Server.Type)
	}

	return srv, err
}
