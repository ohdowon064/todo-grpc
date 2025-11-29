# todo-grpc

간단한 Todo 리스트를 gRPC/Connect 기반 API로 제공하는 Go 프로젝트입니다.

## 주요 기능

- Todo 추가: 제목을 입력하여 새로운 할 일 생성
- Todo 목록 조회: 등록된 모든 할 일 반환

## 기술 스택

- Go 1.25.1
- [ConnectRPC](https://connectrpc.com/)
- GORM (SQLite 사용)
- Protocol Buffers (proto3)
- gRPC/Connect/HTTP2 지원

## 폴더 구조

```
.
├── main.go                # 서버 진입점
├── main_test.go           # 단위 테스트
├── todo.proto             # 프로토콜 버퍼 정의
├── gen/                   # 코드 생성 결과물
│   └── todo/
│       └── v1/
│           ├── todo.pb.go
│           └── todov1connect/
│               └── todo.connect.go
├── Makefile               # 빌드 및 실행 명령어
├── go.mod                 # Go 모듈 설정
└── .gitignore
```

## 실행 방법

1. **의존성 설치**

   ```sh
   go mod tidy
   ```

2. **서버 실행**

   ```sh
   make run
   ```
   또는
   ```sh
   go run main.go
   ```

   서버가 `http://localhost:8080`에서 시작됩니다.

3. **테스트 실행**

   ```sh
   make test
   ```
   또는
   ```sh
   go test -v -cover ./...
   ```

## 실행 및 관리 명령어 (Makefile)

| 명령어         | 설명                                      |
| -------------- | ----------------------------------------- |
| `make run`       | 서버를 로컬에서 실행합니다.                |
| `make test`      | 테스트를 실행합니다. 커버리지 포함         |
| `make add_todo title='할 일'` | 새로운 Todo를 추가합니다.         |
| `make list_todos`| 모든 Todo 항목을 나열합니다.               |
| `make help`      | 사용 가능한 명령어 목록을 표시합니다.      |
| `make gen`      | Protobuf 명세에 대한 go 코드를 생성합니다.      |

### 예시

```sh
make run
make test
make add_todo title="청소하기"
make list_todos
```

## API 명세

### TodoService

#### AddTodo

- **Request:**  
  ```json
  { "title": "할 일 제목" }
  ```
- **Response:**  
  ```json
  { "id": 1 }
  ```

#### ListTodos

- **Request:**  
  ```json
  {}
  ```
- **Response:**  
  ```json
  {
    "todos": [
      { "id": 1, "title": "할 일 제목", "done": false }
    ]
  }
  ```

## 코드 생성 (proto → Go)

프로토콜 버퍼 파일(`todo.proto`)을 수정한 경우 아래 명령어로 코드를 재생성할 수 있습니다.

```sh
make gen 
```

또는 

```sh
protoc --go_out=gen/todo/v1 \
    --go_opt=paths=source_relative \
    --connect-go_out=gen/todo/v1 \
    --connect-go_opt=paths=source_relative \
    todo.proto
```

### 예시

```sh
make gen
```

## 참고

- [ConnectRPC 공식 문서](https://connectrpc.com/docs/go/)
- [GORM 공식 문서](https://gorm.io/docs/)
- [Protocol Buffers 공식 문서](https://protobuf.dev/)

---

문의 및 개선 제안은 이슈로 남겨주세요.