import 'task_list_response.dart';

abstract class TaskApi {
  Future<TaskListResponse> getTasks();
  Future<TaskListResponse> addTask(String task);
  Future<TaskListResponse> removeTask(String task);
}
