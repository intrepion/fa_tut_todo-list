# Finish

Start the API server from the repository root:

```bash
just run
```

This API is configured to accept browser requests from `http://localhost:25616` and to persist tasks in `workspace/data/tasks.db`.

In another terminal, try these requests:

```bash
curl "http://localhost:25664/api/tasks"
curl -X POST "http://localhost:25664/api/tasks" \
  -H "Content-Type: application/json" \
  -d '{"text":"Buy milk"}'
curl "http://localhost:25664/api/tasks/1"
curl -i -X DELETE "http://localhost:25664/api/tasks/1"
```

With a fresh database, the first `GET` should return an empty list:

```json
{"tasks":[]}
```

The `POST` should return the created task resource:

```json
{"id":1,"text":"Buy milk"}
```

The next `GET /api/tasks/1` should return the same task resource. The `DELETE` should return `204 No Content`.

If you already have rows in the database, SQLite may assign a different numeric id than `1`.
