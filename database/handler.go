package database

import (
	"fmt"

	"github.com/stphivos/todo-api-go-grpc/database/mongo"
	"github.com/stphivos/todo-api-go-grpc/models"
)

// Handler ...
type Handler interface {
	GetTodos() ([]models.Todo, error)
}

// Create ...
func Create(config *models.Config) (Handler, error) {
	var db Handler
	var err error

	switch config.Database.Type {
	case "mongo":
		db, err = mongo.NewHandler(config)
	default:
		err = fmt.Errorf("Database type %v is not supported", config.Database.Type)
	}

	return db, err
}
