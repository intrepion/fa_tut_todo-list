# Code

### 1. Red: Load The Current Task List

Create the first code test file:

```bash
touch workspace/test/code/task_list_controller_test.dart
just format
git add --all
git commit --message 'touch workspace/test/code/task_list_controller_test.dart'
```

Put this exact content in `workspace/test/code/task_list_controller_test.dart`:

```dart
import 'package:mocktail/mocktail.dart';
import 'package:todo_list/code/task_list_controller.dart';
import 'package:todo_list/contracts/task_api.dart';
import 'package:todo_list/contracts/task_list_response.dart';
import 'package:test/test.dart';

class MockTaskApi extends Mock implements TaskApi {}

void main() {
  test('loads the current task list', () async {
    final api = MockTaskApi();
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees', 'Buy milk'],
        lines: ['Learn how to invert binary trees', 'Buy milk'],
      ),
    );

    final result = await loadTasks(api);

    expect(result.tasks, ['Learn how to invert binary trees', 'Buy milk']);
    expect(result.lines, ['Learn how to invert binary trees', 'Buy milk']);
    expect(result.errorMessage, isNull);
  });
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "1. Red: Load The Current Task List"
```

### 2. Green: Load The Current Task List

Create the first production file:

```bash
touch workspace/lib/code/task_list_controller.dart
just format
git add --all
git commit --message 'touch workspace/lib/code/task_list_controller.dart'
```

Put this exact content in `workspace/lib/code/task_list_controller.dart`:

```dart
import '../contracts/task_api.dart';
import '../contracts/task_list_response.dart';

class TaskListViewModel {
  final List<String> tasks;
  final List<String> lines;
  final String? errorMessage;

  const TaskListViewModel({
    required this.tasks,
    required this.lines,
    this.errorMessage,
  });

  factory TaskListViewModel.fromResponse(TaskListResponse response) {
    return TaskListViewModel(
      tasks: response.tasks,
      lines: response.lines,
    );
  }
}

Future<TaskListViewModel> loadTasks(TaskApi api) async {
  final response = await api.getTasks();
  return TaskListViewModel.fromResponse(response);
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "2. Green: Load The Current Task List"
```

### 3. Red: Trim Submitted Tasks Before Calling The API

Replace `workspace/test/code/task_list_controller_test.dart` with:

```dart
import 'package:mocktail/mocktail.dart';
import 'package:todo_list/code/task_list_controller.dart';
import 'package:todo_list/contracts/task_api.dart';
import 'package:todo_list/contracts/task_list_response.dart';
import 'package:test/test.dart';

class MockTaskApi extends Mock implements TaskApi {}

void main() {
  test('loads the current task list', () async {
    final api = MockTaskApi();
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees', 'Buy milk'],
        lines: ['Learn how to invert binary trees', 'Buy milk'],
      ),
    );

    final result = await loadTasks(api);

    expect(result.tasks, ['Learn how to invert binary trees', 'Buy milk']);
    expect(result.lines, ['Learn how to invert binary trees', 'Buy milk']);
    expect(result.errorMessage, isNull);
  });

  test('trims submitted tasks before calling the API', () async {
    final api = MockTaskApi();
    when(() => api.addTask('Buy milk')).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees', 'Buy milk'],
        lines: ['Learn how to invert binary trees', 'Buy milk'],
      ),
    );

    final result = await addTask('  Buy milk  ', api);

    expect(result.tasks, ['Learn how to invert binary trees', 'Buy milk']);
    verify(() => api.addTask('Buy milk')).called(1);
  });
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "3. Red: Trim Submitted Tasks Before Calling The API"
```

### 4. Green: Trim Submitted Tasks Before Calling The API

Replace `workspace/lib/code/task_list_controller.dart` with:

```dart
import '../contracts/task_api.dart';
import '../contracts/task_list_response.dart';

class TaskListViewModel {
  final List<String> tasks;
  final List<String> lines;
  final String? errorMessage;

  const TaskListViewModel({
    required this.tasks,
    required this.lines,
    this.errorMessage,
  });

  factory TaskListViewModel.fromResponse(TaskListResponse response) {
    return TaskListViewModel(
      tasks: response.tasks,
      lines: response.lines,
    );
  }
}

Future<TaskListViewModel> loadTasks(TaskApi api) async {
  final response = await api.getTasks();
  return TaskListViewModel.fromResponse(response);
}

Future<TaskListViewModel> addTask(String task, TaskApi api) async {
  final response = await api.addTask(task.trim());
  return TaskListViewModel.fromResponse(response);
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "4. Green: Trim Submitted Tasks Before Calling The API"
```

### 5. Red: Reject Blank Tasks And Handle API Errors

Replace `workspace/test/code/task_list_controller_test.dart` with:

```dart
import 'package:mocktail/mocktail.dart';
import 'package:todo_list/code/task_list_controller.dart';
import 'package:todo_list/contracts/task_api.dart';
import 'package:todo_list/contracts/task_list_response.dart';
import 'package:test/test.dart';

class MockTaskApi extends Mock implements TaskApi {}

void main() {
  test('loads the current task list', () async {
    final api = MockTaskApi();
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees', 'Buy milk'],
        lines: ['Learn how to invert binary trees', 'Buy milk'],
      ),
    );

    final result = await loadTasks(api);

    expect(result.tasks, ['Learn how to invert binary trees', 'Buy milk']);
    expect(result.lines, ['Learn how to invert binary trees', 'Buy milk']);
    expect(result.errorMessage, isNull);
  });

  test('trims submitted tasks before calling the API', () async {
    final api = MockTaskApi();
    when(() => api.addTask('Buy milk')).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees', 'Buy milk'],
        lines: ['Learn how to invert binary trees', 'Buy milk'],
      ),
    );

    final result = await addTask('  Buy milk  ', api);

    expect(result.tasks, ['Learn how to invert binary trees', 'Buy milk']);
    verify(() => api.addTask('Buy milk')).called(1);
  });

  test('rejects blank submitted tasks without calling the API', () async {
    final api = MockTaskApi();

    final result = await addTask('   ', api);

    expect(result.errorMessage, 'Task must not be blank.');
    verifyNever(() => api.addTask(any()));
  });

  test('returns a friendly message when the task API is unavailable', () async {
    final api = MockTaskApi();
    when(() => api.getTasks()).thenThrow(Exception('boom'));

    final result = await loadTasks(api);

    expect(result.errorMessage, 'Sorry, the task API is unavailable right now.');
  });

  test('removes the chosen task through the API', () async {
    final api = MockTaskApi();
    when(() => api.removeTask('Buy milk')).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees'],
        lines: ['Learn how to invert binary trees'],
      ),
    );

    final result = await removeTask('Buy milk', api);

    expect(result.tasks, ['Learn how to invert binary trees']);
    verify(() => api.removeTask('Buy milk')).called(1);
  });
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "5. Red: Reject Blank Tasks And Handle API Errors"
```

### 6. Green: Reject Blank Tasks And Handle API Errors

Replace `workspace/lib/code/task_list_controller.dart` with:

```dart
import '../contracts/task_api.dart';
import '../contracts/task_list_response.dart';

class TaskListViewModel {
  final List<String> tasks;
  final List<String> lines;
  final String? errorMessage;

  const TaskListViewModel({
    required this.tasks,
    required this.lines,
    this.errorMessage,
  });

  factory TaskListViewModel.fromResponse(TaskListResponse response) {
    return TaskListViewModel(
      tasks: response.tasks,
      lines: response.lines,
    );
  }
}

Future<TaskListViewModel> loadTasks(TaskApi api) async {
  try {
    final response = await api.getTasks();
    return TaskListViewModel.fromResponse(response);
  } catch (_) {
    return const TaskListViewModel(
      tasks: [],
      lines: [],
      errorMessage: 'Sorry, the task API is unavailable right now.',
    );
  }
}

Future<TaskListViewModel> addTask(String task, TaskApi api) async {
  final trimmed = task.trim();
  if (trimmed.isEmpty) {
    return const TaskListViewModel(
      tasks: [],
      lines: [],
      errorMessage: 'Task must not be blank.',
    );
  }

  try {
    final response = await api.addTask(trimmed);
    return TaskListViewModel.fromResponse(response);
  } catch (_) {
    return const TaskListViewModel(
      tasks: [],
      lines: [],
      errorMessage: 'Sorry, the task API is unavailable right now.',
    );
  }
}

Future<TaskListViewModel> removeTask(String task, TaskApi api) async {
  try {
    final response = await api.removeTask(task);
    return TaskListViewModel.fromResponse(response);
  } catch (_) {
    return const TaskListViewModel(
      tasks: [],
      lines: [],
      errorMessage: 'Sorry, the task API is unavailable right now.',
    );
  }
}
```

Run:

```bash
just format
just check-all
git add --all
git commit --message "6. Green: Reject Blank Tasks And Handle API Errors"
```
