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

(cd workspace && go get github.com/labstack/echo/v4)
just format
git add --all
git commit --message "(cd workspace && go get github.com/labstack/echo/v4)"

(cd workspace && go get github.com/labstack/echo/v4/middleware)
just format
git add --all
git commit --message "(cd workspace && go get github.com/labstack/echo/v4/middleware)"

(cd workspace && go get modernc.org/sqlite)
just format
git add --all
git commit --message "(cd workspace && go get modernc.org/sqlite)"

(cd workspace && go get github.com/stretchr/testify/assert github.com/stretchr/testify/mock)
just format
git add --all
git commit --message "(cd workspace && go get github.com/stretchr/testify/assert github.com/stretchr/testify/mock)"

mkdir -p workspace/data
```

This gives you:

- a root-level `.gitignore` for operating-system noise and editor leftovers
- a `workspace/.gitignore` for standard Go build output and local tooling files
- a `workspace/data/tasks.db` file path for the durable SQLite task store used by the REST API adapter

When the full workspace is finished, it should contain these files:

```text
workspace/
  .gitignore
  go.mod
  go.sum
  cmd/
    server/
      main.go
  data/
    tasks.db
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
        sqlite_task_store.go
        sqlite_task_store_test.go
```
