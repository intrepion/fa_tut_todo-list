import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:todo_list/adapter/http_task_api.dart';
import 'package:todo_list/contracts/task.dart';
import 'package:todo_list/contracts/task_list_response.dart';
import 'package:test/test.dart';

void main() {
  test('loads the current tasks from the canonical endpoint', () async {
    late Uri requestedUri;
    final client = MockClient((request) async {
      requestedUri = request.url;
      return http.Response(
        '{"tasks":[{"id":"11111111-1111-1111-1111-111111111111","text":"Learn how to invert binary trees"},{"id":"22222222-2222-2222-2222-222222222222","text":"Buy milk"}]}',
        200,
        headers: {'content-type': 'application/json'},
      );
    });

    final api = HttpTaskApi(baseUrl: 'http://localhost:25664', client: client);
    final result = await api.getTasks();

    expect(requestedUri.toString(), 'http://localhost:25664/api/tasks');
    expect(result.tasks.map((task) => task.text).toList(), [
      'Learn how to invert binary trees',
      'Buy milk',
    ]);
  });

  test('posts a new task resource to the canonical endpoint', () async {
    late String requestBody;
    final client = MockClient((request) async {
      requestBody = request.body;
      return http.Response(
        '{"id":"22222222-2222-2222-2222-222222222222","text":"Buy milk"}',
        201,
        headers: {'content-type': 'application/json'},
      );
    });

    final api = HttpTaskApi(baseUrl: 'http://localhost:25664', client: client);
    final createdTask = await api.createTask('Buy milk');

    expect(requestBody, '{"text":"Buy milk"}');
    expect(createdTask, isA<Task>());
    expect(createdTask.id, '22222222-2222-2222-2222-222222222222');
    expect(createdTask.text, 'Buy milk');
  });

  test('deletes a task through the canonical endpoint', () async {
    late Uri requestedUri;
    final client = MockClient((request) async {
      requestedUri = request.url;
      return http.Response('', 204);
    });

    final api = HttpTaskApi(baseUrl: 'http://localhost:25664', client: client);
    await api.deleteTask('77777777-7777-7777-7777-777777777777');

    expect(
      requestedUri.toString(),
      'http://localhost:25664/api/tasks/77777777-7777-7777-7777-777777777777',
    );
  });
}
