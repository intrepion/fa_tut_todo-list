# Contracts

Create the shared contract files:

```bash
mkdir -p workspace/lib/contracts
touch workspace/lib/contracts/task_list_response.dart
just format
git add --all
git commit --message 'touch workspace/lib/contracts/task_list_response.dart'
touch workspace/lib/contracts/task_api.dart
just format
git add --all
git commit --message 'touch workspace/lib/contracts/task_api.dart'
```

Put this exact content in `workspace/lib/contracts/task_list_response.dart`:

```dart
class TaskListResponse {
  final List<String> tasks;
  final List<String> lines;

  const TaskListResponse({required this.tasks, required this.lines});

  factory TaskListResponse.fromJson(Map<String, dynamic> json) {
    return TaskListResponse(
      tasks: List<String>.from(json['tasks'] as List<dynamic>),
      lines: List<String>.from(json['lines'] as List<dynamic>),
    );
  }
}
```

Put this exact content in `workspace/lib/contracts/task_api.dart`:

```dart
import 'task_list_response.dart';

abstract class TaskApi {
  Future<TaskListResponse> getTasks();
  Future<TaskListResponse> addTask(String task);
  Future<TaskListResponse> removeTask(String task);
}
```

Do not add tests here. Keep this layer limited to interfaces and small shared types.

Then run:

```bash
just format
just check-all
git add --all
git commit --message "Define todo-list Flutter contracts"
```
