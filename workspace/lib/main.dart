import 'package:flutter/material.dart';

import 'adapter/http_task_api.dart';
import 'adapter/todo_list_page.dart';

const apiBaseUrl = String.fromEnvironment(
  'API_BASE_URL',
  defaultValue: 'http://localhost:25664',
);

void main() {
  runApp(const TodoListApp());
}

class TodoListApp extends StatelessWidget {
  const TodoListApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Todo List',
      home: TodoListPage(api: HttpTaskApi(baseUrl: apiBaseUrl)),
    );
  }
}
