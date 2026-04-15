set shell := ["bash", "-eu", "-c"]

default:
	@just --list

workspace := "workspace"

restore:
	(cd {{workspace}} && go mod download)

format:
	find {{workspace}} -name '*.go' -exec gofmt -w {} +

check-formatting:
	test -z "$(find {{workspace}} -name '*.go' -exec gofmt -l {} +)"

check-tests:
	(cd {{workspace}} && go test ./...)

run:
	(cd {{workspace}} && go run ./cmd/server)

check-all:
	just check-formatting
	just check-tests
