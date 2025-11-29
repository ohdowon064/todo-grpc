package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	todov1 "todo-grpc/gen/todo/v1"
	"todo-grpc/gen/todo/v1/todov1connect"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TodoModel struct {
	gorm.Model
	Title string `gorm:"not null"`
	Done  bool   `gorm:"default:false"`
}

func (TodoModel) TableName() string {
	return "todos"
}

type TodoServer struct {
	db *gorm.DB
	todov1connect.UnimplementedTodoServiceHandler
}

func (s *TodoServer) AddTodo(
	ctx context.Context,
	req *connect.Request[todov1.AddTodoRequest],
) (*connect.Response[todov1.AddTodoResponse], error) {
	log.Printf("Request AddTodo: %s", req.Msg.Title)

	newTodo := TodoModel{Title: req.Msg.Title}
	if result := s.db.Create(&newTodo); result.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, result.Error)
	}

	return connect.NewResponse(
		&todov1.AddTodoResponse{
			Id: int64(newTodo.ID),
		}), nil
}

func (s *TodoServer) ListTodos(
	ctx context.Context,
	req *connect.Request[todov1.ListTodosRequest],
) (*connect.Response[todov1.ListTodosResponse], error) {
	log.Printf("Request ListTodos %s", req.Msg)

	var models []TodoModel
	if result := s.db.Find(&models); result.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, result.Error)
	}

	var protos []*todov1.Todo
	for _, model := range models {
		protos = append(protos, &todov1.Todo{
			Id:    int64(model.ID),
			Title: model.Title,
			Done:  model.Done,
		})
	}

	return connect.NewResponse(&todov1.ListTodosResponse{
		Todos: protos,
	}), nil
}

func main() {
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(&TodoModel{}); err != nil {
		log.Fatal(err)
	}

	todoServer := &TodoServer{db: db}
	mux := http.NewServeMux()

	path, handler := todov1connect.NewTodoServiceHandler(todoServer)
	mux.Handle(path, handler)

	fmt.Println("Starting server on http://localhost:8080")

	err = http.ListenAndServe(
		":8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatal(err)
	}
}
