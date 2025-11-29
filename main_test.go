package main

import (
	"context"
	"testing"
	todov1 "todo-grpc/gen/todo/v1"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestServer(t *testing.T) *TodoServer {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("DB 연결 실패: %v", err)
	}

	err = db.AutoMigrate(&TodoModel{})
	if err != nil {
		t.Fatalf("마이그레이션 실패: %v", err)
	}

	return &TodoServer{db: db}
}

func TestTodoLifecycle(t *testing.T) {
	server := setupTestServer(t)
	ctx := context.Background()

	t.Run("Add Todo", func(t *testing.T) {
		req := connect.NewRequest(&todov1.AddTodoRequest{
			Title: "CI 테스트 통과하기",
		})

		res, err := server.AddTodo(ctx, req)
		assert.NoError(t, err)
		assert.NotZero(t, res.Msg.Id)
		assert.Equal(t, int64(1), res.Msg.Id)
	})

	t.Run("List Todos", func(t *testing.T) {
		req := connect.NewRequest(&todov1.ListTodosRequest{})
		res, err := server.ListTodos(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, res.Msg.Todos, 1)
		assert.Equal(t, "CI 테스트 통과하기", res.Msg.Todos[0].Title)
		assert.False(t, res.Msg.Todos[0].Done)
	})
}
