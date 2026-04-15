import 'task.dart';
import 'task_list_response.dart';

abstract class TaskApi {
  Future<TaskListResponse> getTasks();
  Future<Task> createTask(String text);
  Future<void> deleteTask(int taskId);
}
