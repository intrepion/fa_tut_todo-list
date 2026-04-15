set shell := ["bash", "-eu", "-c"]

default:
	@just --list

workspace := "workspace"
database_url := "postgres://postgres@localhost:5432/todo_list?sslmode=disable"

restore:
	(cd {{workspace}} && go mod download)

format:
	find {{workspace}} -name '*.go' -exec gofmt -w {} +

check-formatting:
	test -z "$(find {{workspace}} -name '*.go' -exec gofmt -l {} +)"

check-tests:
	(cd {{workspace}} && go test ./...)

run:
	(cd {{workspace}} && TODO_LIST_DATABASE_URL={{database_url}} go run ./cmd/server)

check-all:
	just check-formatting
	just check-tests
