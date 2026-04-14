# Contracts

Create the shared contract files:

```bash
mkdir -p workspace/lib/contracts
touch workspace/lib/contracts/task.dart
just format
git add --all
git commit --message 'touch workspace/lib/contracts/task.dart'
touch workspace/lib/contracts/task_list_response.dart
just format
git add --all
git commit --message 'touch workspace/lib/contracts/task_list_response.dart'
touch workspace/lib/contracts/task_api.dart
just format
git add --all
git commit --message 'touch workspace/lib/contracts/task_api.dart'
```

Put this exact content in `workspace/lib/contracts/task.dart`:

```dart
class Task {
  final int id;
  final String text;

  const Task({required this.id, required this.text});

  factory Task.fromJson(Map<String, dynamic> json) {
    return Task(
      id: json['id'] as int,
      text: json['text'] as String,
    );
  }
}
```

Put this exact content in `workspace/lib/contracts/task_list_response.dart`:

```dart
import 'task.dart';

class TaskListResponse {
  final List<Task> tasks;

  const TaskListResponse({required this.tasks});

  factory TaskListResponse.fromJson(Map<String, dynamic> json) {
    return TaskListResponse(
      tasks: (json['tasks'] as List<dynamic>)
          .map((taskJson) => Task.fromJson(taskJson as Map<String, dynamic>))
          .toList(),
    );
  }
}
```

Put this exact content in `workspace/lib/contracts/task_api.dart`:

```dart
import 'task.dart';
import 'task_list_response.dart';

abstract class TaskApi {
  Future<TaskListResponse> getTasks();
  Future<Task> createTask(String text);
  Future<void> deleteTask(int taskId);
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
