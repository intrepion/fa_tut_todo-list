import 'dart:convert';

import 'package:http/http.dart' as http;

import '../contracts/task_api.dart';
import '../contracts/task_list_response.dart';

class HttpTaskApi implements TaskApi {
  final String baseUrl;
  final http.Client client;

  HttpTaskApi({required this.baseUrl, http.Client? client})
    : client = client ?? http.Client();

  @override
  Future<TaskListResponse> getTasks() async {
    final response = await client.get(Uri.parse('$baseUrl/api/tasks'));
    return TaskListResponse.fromJson(
      jsonDecode(response.body) as Map<String, dynamic>,
    );
  }

  @override
  Future<TaskListResponse> addTask(String task) async {
    final response = await client.post(
      Uri.parse('$baseUrl/api/tasks'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'task': task}),
    );
    return TaskListResponse.fromJson(
      jsonDecode(response.body) as Map<String, dynamic>,
    );
  }

  @override
  Future<TaskListResponse> removeTask(String task) async {
    final uri = Uri.parse(
      '$baseUrl/api/tasks',
    ).replace(queryParameters: {'task': task});
    final response = await client.delete(uri);
    return TaskListResponse.fromJson(
      jsonDecode(response.body) as Map<String, dynamic>,
    );
  }
}
