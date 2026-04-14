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
    return TaskListViewModel(tasks: response.tasks, lines: response.lines);
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
