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
  final response = await api.getTasks();
  return TaskListViewModel.fromResponse(response);
}
