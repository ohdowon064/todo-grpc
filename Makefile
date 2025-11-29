default: help

.PHONY: help
help: ## 도움말 메시지를 표시합니다.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## 애플리케이션을 로컬에서 실행합니다.
	go run main.go

.PHONY: test
test: ## 테스트를 실행합니다. 커버리지 포함
	go test -v -cover ./...

.PHONY: gen
gen: ## Protobuf 명세에 대한 go 코드를 생성합니다.
	@mkdir -p gen/todo/v1
	@protoc --go_out=gen/todo/v1 \
	--go_opt=paths=source_relative \
	--connect-go_out=gen/todo/v1 \
	--connect-go_opt=paths=source_relative \
	todo.proto

.PHONY: add_todo
add_todo: ## 새로운 Todo를 추가합니다. make add_todo title='할 일 제목'
	@if [ -z "$(title)" ]; then \
		echo "사용법: make add_todo title='할 일 제목'"; \
		exit 1; \
	fi
	@echo "Todo를 추가합니다. 제목: $(title)"
	@curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"title": "$(title)"}' \
	http://localhost:8080/todo.v1.TodoService/AddTodo

.PHONY: list_todos
list_todos: ## 모든 Todo 항목을 나열합니다.
	@echo "모든 Todo 항목을 나열합니다..."
	@curl -X POST \
	-H "Content-Type: application/json" \
	-d '{}' \
	http://localhost:8080/todo.v1.TodoService/ListTodos