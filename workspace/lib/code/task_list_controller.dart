import '../contracts/task.dart';
import '../contracts/task_api.dart';
import '../contracts/task_list_response.dart';

class TaskListViewModel {
  final List<Task> tasks;
  final String? errorMessage;

  const TaskListViewModel({required this.tasks, this.errorMessage});

  factory TaskListViewModel.fromResponse(TaskListResponse response) {
    return TaskListViewModel(tasks: response.tasks);
  }
}

Future<TaskListViewModel> loadTasks(TaskApi api) async {
  try {
    final response = await api.getTasks();
    return TaskListViewModel.fromResponse(response);
  } catch (_) {
    return const TaskListViewModel(
      tasks: [],
      errorMessage: 'Sorry, the task API is unavailable right now.',
    );
  }
}

Future<TaskListViewModel> addTask(String task, TaskApi api) async {
  final trimmed = task.trim();
  if (trimmed.isEmpty) {
    return const TaskListViewModel(
      tasks: [],
      errorMessage: 'Task must not be blank.',
    );
  }

  try {
    await api.createTask(trimmed);
    final response = await api.getTasks();
    return TaskListViewModel.fromResponse(response);
  } catch (_) {
    return const TaskListViewModel(
      tasks: [],
      errorMessage: 'Sorry, the task API is unavailable right now.',
    );
  }
}

Future<TaskListViewModel> removeTask(int taskId, TaskApi api) async {
  try {
    await api.deleteTask(taskId);
    final response = await api.getTasks();
    return TaskListViewModel.fromResponse(response);
  } catch (_) {
    return const TaskListViewModel(
      tasks: [],
      errorMessage: 'Sorry, the task API is unavailable right now.',
    );
  }
}
