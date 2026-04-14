# Finish

Start the API server from the repository root:

```bash
just run
```

This API is configured to accept browser requests from `http://localhost:25616` and to persist tasks in `workspace/data/tasks.json`.

In another terminal, try these requests:

```bash
curl "http://localhost:25664/api/tasks"
curl -X POST "http://localhost:25664/api/tasks" \
  -H "Content-Type: application/json" \
  -d '{"task":"Buy milk"}'
curl -X DELETE "http://localhost:25664/api/tasks?task=Buy%20milk"
```

The initial `GET` should return an empty task list:

```json
{"tasks":[],"lines":[]}
```

After the `POST`, you should get:

```json
{"tasks":["Buy milk"],"lines":["Buy milk"]}
```

After the `DELETE`, the task list should be empty again.

If you later build a browser client at `http://localhost:25616`, it can call this API without additional CORS setup.
