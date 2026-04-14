# Adapter

### 1. Red: Add The HTTP Task API Test

Create the HTTP adapter test file:

```bash
touch workspace/test/adapter/http_task_api_test.dart
just format
git add --all
git commit --message 'touch workspace/test/adapter/http_task_api_test.dart'
```

Put this exact content in `workspace/test/adapter/http_task_api_test.dart`:

```dart
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:todo_list/adapter/http_task_api.dart';
import 'package:test/test.dart';

void main() {
  test('loads the current tasks from the canonical endpoint', () async {
    late Uri requestedUri;
    final client = MockClient((request) async {
      requestedUri = request.url;
      return http.Response(
        '{"tasks":["Learn how to invert binary trees","Buy milk"],"lines":["Learn how to invert binary trees","Buy milk"]}',
        200,
        headers: {'content-type': 'application/json'},
      );
    });

    final api = HttpTaskApi(
      baseUrl: 'http://localhost:25664',
      client: client,
    );
    final result = await api.getTasks();

    expect(requestedUri.toString(), 'http://localhost:25664/api/tasks');
    expect(result.tasks, ['Learn how to invert binary trees', 'Buy milk']);
  });

  test('posts a new task to the canonical endpoint', () async {
    late String requestBody;
    final client = MockClient((request) async {
      requestBody = request.body;
      return http.Response(
        '{"tasks":["Learn how to invert binary trees","Buy milk"],"lines":["Learn how to invert binary trees","Buy milk"]}',
        200,
        headers: {'content-type': 'application/json'},
      );
    });

    final api = HttpTaskApi(
      baseUrl: 'http://localhost:25664',
      client: client,
    );
    await api.addTask('Buy milk');

    expect(requestBody, '{"task":"Buy milk"}');
  });

  test('deletes a task through the canonical endpoint', () async {
    late Uri requestedUri;
    final client = MockClient((request) async {
      requestedUri = request.url;
      return http.Response(
        '{"tasks":["Learn how to invert binary trees"],"lines":["Learn how to invert binary trees"]}',
        200,
        headers: {'content-type': 'application/json'},
      );
    });

    final api = HttpTaskApi(
      baseUrl: 'http://localhost:25664',
      client: client,
    );
    await api.removeTask('Buy milk');

    expect(
      requestedUri.toString(),
      'http://localhost:25664/api/tasks?task=Buy%20milk',
    );
  });
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "1. Red: Add The HTTP Task API Test"
```

### 2. Green: Call The Canonical Task API Endpoints

Create the HTTP adapter production file:

```bash
touch workspace/lib/adapter/http_task_api.dart
just format
git add --all
git commit --message 'touch workspace/lib/adapter/http_task_api.dart'
```

Put this exact content in `workspace/lib/adapter/http_task_api.dart`:

```dart
import 'dart:convert';

import 'package:http/http.dart' as http;

import '../contracts/task_api.dart';
import '../contracts/task_list_response.dart';

class HttpTaskApi implements TaskApi {
  final String baseUrl;
  final http.Client client;

  HttpTaskApi({required this.baseUrl, http.Client? client})
    : client = client ?? http.Client();

  @override
  Future<TaskListResponse> getTasks() async {
    final response = await client.get(Uri.parse('$baseUrl/api/tasks'));
    return TaskListResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
  }

  @override
  Future<TaskListResponse> addTask(String task) async {
    final response = await client.post(
      Uri.parse('$baseUrl/api/tasks'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'task': task}),
    );
    return TaskListResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
  }

  @override
  Future<TaskListResponse> removeTask(String task) async {
    final response = await client.delete(
      Uri.parse('$baseUrl/api/tasks?task=${Uri.encodeQueryComponent(task)}'),
    );
    return TaskListResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
  }
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "2. Green: Call The Canonical Task API Endpoints"
```

### 3. Red: Add The Todo List Page Widget Test

Create the widget test file:

```bash
touch workspace/test/adapter/todo_list_page_test.dart
just format
git add --all
git commit --message 'touch workspace/test/adapter/todo_list_page_test.dart'
```

Put this exact content in `workspace/test/adapter/todo_list_page_test.dart`:

```dart
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:todo_list/adapter/todo_list_page.dart';
import 'package:todo_list/contracts/task_api.dart';
import 'package:todo_list/contracts/task_list_response.dart';

class MockTaskApi extends Mock implements TaskApi {}

void main() {
  testWidgets('loads, adds, and removes tasks', (tester) async {
    final api = MockTaskApi();
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees'],
        lines: ['Learn how to invert binary trees'],
      ),
    );
    when(() => api.addTask('Buy milk')).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees', 'Buy milk'],
        lines: ['Learn how to invert binary trees', 'Buy milk'],
      ),
    );
    when(() => api.removeTask('Buy milk')).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees'],
        lines: ['Learn how to invert binary trees'],
      ),
    );

    await tester.pumpWidget(
      MaterialApp(
        home: TodoListPage(api: api),
      ),
    );
    await tester.pumpAndSettle();

    expect(find.text('Learn how to invert binary trees'), findsOneWidget);

    await tester.enterText(find.byType(TextField), 'Buy milk');
    await tester.tap(find.text('Add task'));
    await tester.pumpAndSettle();

    expect(find.text('Buy milk'), findsOneWidget);

    await tester.tap(find.byKey(const Key('remove-Buy milk')));
    await tester.pumpAndSettle();

    expect(find.text('Buy milk'), findsNothing);
  });
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "3. Red: Add The Todo List Page Widget Test"
```

### 4. Green: Build The Todo List Page

Create the page production file:

```bash
touch workspace/lib/adapter/todo_list_page.dart
just format
git add --all
git commit --message 'touch workspace/lib/adapter/todo_list_page.dart'
```

Put this exact content in `workspace/lib/adapter/todo_list_page.dart`:

```dart
import 'package:flutter/material.dart';

import '../code/task_list_controller.dart';
import '../contracts/task_api.dart';

class TodoListPage extends StatefulWidget {
  final TaskApi api;

  const TodoListPage({super.key, required this.api});

  @override
  State<TodoListPage> createState() => _TodoListPageState();
}

class _TodoListPageState extends State<TodoListPage> {
  final TextEditingController _taskController = TextEditingController();
  TaskListViewModel _viewModel = const TaskListViewModel(tasks: [], lines: []);

  @override
  void initState() {
    super.initState();
    _refreshTasks();
  }

  Future<void> _refreshTasks() async {
    final nextViewModel = await loadTasks(widget.api);
    if (!mounted) {
      return;
    }
    setState(() {
      _viewModel = nextViewModel;
    });
  }

  Future<void> _addTask() async {
    final nextViewModel = await addTask(_taskController.text, widget.api);
    _taskController.clear();
    if (!mounted) {
      return;
    }
    setState(() {
      _viewModel = nextViewModel;
    });
  }

  Future<void> _removeTask(String task) async {
    final nextViewModel = await removeTask(task, widget.api);
    if (!mounted) {
      return;
    }
    setState(() {
      _viewModel = nextViewModel;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Todo List')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _taskController,
              decoration: const InputDecoration(labelText: 'Task'),
            ),
            const SizedBox(height: 12),
            ElevatedButton(
              onPressed: _addTask,
              child: const Text('Add task'),
            ),
            if (_viewModel.errorMessage != null) ...[
              const SizedBox(height: 12),
              Text(_viewModel.errorMessage!),
            ],
            const SizedBox(height: 12),
            Expanded(
              child: ListView(
                children: _viewModel.tasks
                    .map(
                      (task) => ListTile(
                        title: Text(task),
                        trailing: IconButton(
                          key: Key('remove-$task'),
                          onPressed: () => _removeTask(task),
                          icon: const Icon(Icons.delete_outline),
                        ),
                      ),
                    )
                    .toList(),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "4. Green: Build The Todo List Page"
```

### 5. Green: Wire The Real Application

Replace `workspace/lib/main.dart` with:

```dart
import 'package:flutter/material.dart';

import 'adapter/http_task_api.dart';
import 'adapter/todo_list_page.dart';

const apiBaseUrl = String.fromEnvironment(
  'API_BASE_URL',
  defaultValue: 'http://localhost:25664',
);

void main() {
  runApp(const TodoListApp());
}

class TodoListApp extends StatelessWidget {
  const TodoListApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Todo List',
      home: TodoListPage(
        api: HttpTaskApi(baseUrl: apiBaseUrl),
      ),
    );
  }
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "5. Green: Wire The Real Application"
```
