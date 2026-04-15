import 'dart:convert';

import 'package:http/http.dart' as http;

import '../contracts/task_api.dart';
import '../contracts/task.dart';
import '../contracts/task_list_response.dart';

class HttpTaskApi implements TaskApi {
  final String baseUrl;
  final http.Client client;

  HttpTaskApi({required this.baseUrl, http.Client? client})
    : client = client ?? http.Client();

  @override
  Future<TaskListResponse> getTasks() async {
    final response = await client.get(Uri.parse('$baseUrl/api/tasks'));
    if (response.statusCode != 200) {
      throw Exception('failed to load tasks');
    }
    return TaskListResponse.fromJson(
      jsonDecode(response.body) as Map<String, dynamic>,
    );
  }

  @override
  Future<Task> createTask(String text) async {
    final response = await client.post(
      Uri.parse('$baseUrl/api/tasks'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'text': text}),
    );
    if (response.statusCode != 201) {
      throw Exception('failed to create task');
    }
    return Task.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
  }

  @override
  Future<void> deleteTask(int taskId) async {
    final response = await client.delete(
      Uri.parse('$baseUrl/api/tasks/$taskId'),
    );
    if (response.statusCode != 204) {
      throw Exception('failed to delete task');
    }
  }
}
