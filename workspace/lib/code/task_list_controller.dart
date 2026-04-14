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
  final response = await api.getTasks();
  return TaskListViewModel.fromResponse(response);
}
