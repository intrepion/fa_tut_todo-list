class Task {
  final int id;
  final String text;

  const Task({required this.id, required this.text});

  factory Task.fromJson(Map<String, dynamic> json) {
    return Task(id: json['id'] as int, text: json['text'] as String);
  }
}
