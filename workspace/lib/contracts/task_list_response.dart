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
