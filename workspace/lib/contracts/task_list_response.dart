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
