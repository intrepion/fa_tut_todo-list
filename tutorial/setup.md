# Setup

Keep the repository root for shared files like `README.md`, `LICENSE`, `.gitignore`, `.github/`, `justfile`, and `tutorial/`.

Put all Go code inside a single `workspace/` folder.

From the repository root, run each setup command and checkpoint it before moving to the next one:

```bash
rm -rf workspace

mkdir -p workspace

curl -L -s https://raw.githubusercontent.com/github/gitignore/refs/heads/main/Go.gitignore > workspace/.gitignore
just format
git add --all
git commit --message "curl -L -s https://raw.githubusercontent.com/github/gitignore/refs/heads/main/Go.gitignore > workspace/.gitignore"

(cd workspace && go mod init github.com/intrepion/fa_tut_todo-list/workspace)
just format
git add --all
git commit --message "(cd workspace && go mod init github.com/intrepion/fa_tut_todo-list/workspace)"

(cd workspace && GOBIN=$(pwd)/bin go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0)
just format
git add --all
git commit --message "(cd workspace && GOBIN=\$(pwd)/bin go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0)"

(cd workspace && go get github.com/labstack/echo/v4)
just format
git add --all
git commit --message "(cd workspace && go get github.com/labstack/echo/v4)"

(cd workspace && go get github.com/labstack/echo/v4/middleware)
just format
git add --all
git commit --message "(cd workspace && go get github.com/labstack/echo/v4/middleware)"

(cd workspace && go get github.com/jackc/pgx/v5/stdlib)
just format
git add --all
git commit --message "(cd workspace && go get github.com/jackc/pgx/v5/stdlib)"

(cd workspace && go get github.com/DATA-DOG/go-sqlmock)
just format
git add --all
git commit --message "(cd workspace && go get github.com/DATA-DOG/go-sqlmock)"

(cd workspace && go get github.com/stretchr/testify/assert github.com/stretchr/testify/mock)
just format
git add --all
git commit --message "(cd workspace && go get github.com/stretchr/testify/assert github.com/stretchr/testify/mock)"
```

This gives you:

- a root-level `.gitignore` for operating-system noise and editor leftovers
- a `workspace/.gitignore` for standard Go build output and local tooling files
- a Postgres-backed REST adapter that reads its connection string from `TODO_LIST_DATABASE_URL`
- a generated root `justfile` that defaults `database_url` to `postgres://postgres@localhost:5432/todo_list?sslmode=disable`
- a repo-local `workspace/bin/sqlc` installation for generating Go query code from `workspace/sqlc.yaml`, `workspace/db/schema.sql`, and `workspace/db/query/tasks.sql`

Before you run the server, create the default tutorial database with:

```bash
createdb --host localhost --username postgres todo_list
```

If that does not match your local Postgres setup, create an equivalent database your user can access and override the generated `database_url` value in the root `justfile`.

When the full workspace is finished, it should contain these files:

```text
workspace/
  .gitignore
  bin/
    sqlc
  go.mod
  go.sum
  sqlc.yaml
  db/
    schema.sql
    query/
      tasks.sql
  cmd/
    server/
      main.go
  internal/
    contracts/
      task_api.go
    code/
      task_service.go
      task_service_test.go
    adapter/
      http/
        task_handler.go
        task_handler_test.go
      storage/
        postgres_task_store.go
        postgres_task_store_test.go
        db/
          ...generated Go files from sqlc...
```
