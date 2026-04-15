set shell := ["bash", "-eu", "-c"]

default:
	@just --list

workspace := "workspace"
database_url := "postgres://postgres@localhost:5432/todo_list?sslmode=disable"

restore:
	(cd {{workspace}} && go mod download)

generate:
	(cd {{workspace}} && sqlc generate)

format:
	find {{workspace}} -name '*.go' -exec gofmt -w {} +

check-formatting:
	test -z "$(find {{workspace}} -name '*.go' -exec gofmt -l {} +)"

check-tests:
	just generate
	(cd {{workspace}} && go test ./...)

run:
	just generate
	(cd {{workspace}} && TODO_LIST_DATABASE_URL={{database_url}} go run ./cmd/server)

check-all:
	just check-formatting
	just check-tests
