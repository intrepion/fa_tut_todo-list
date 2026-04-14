import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:todo_list/adapter/http_task_api.dart';
import 'package:test/test.dart';

void main() {
  test('loads the current tasks from the canonical endpoint', () async {
    late Uri requestedUri;
    final client = MockClient((request) async {
      requestedUri = request.url;
      return http.Response(
        '{"tasks":["Learn how to invert binary trees","Buy milk"],"lines":["Learn how to invert binary trees","Buy milk"]}',
        200,
        headers: {'content-type': 'application/json'},
      );
    });

    final api = HttpTaskApi(baseUrl: 'http://localhost:25664', client: client);
    final result = await api.getTasks();

    expect(requestedUri.toString(), 'http://localhost:25664/api/tasks');
    expect(result.tasks, ['Learn how to invert binary trees', 'Buy milk']);
  });

  test('posts a new task to the canonical endpoint', () async {
    late String requestBody;
    final client = MockClient((request) async {
      requestBody = request.body;
      return http.Response(
        '{"tasks":["Learn how to invert binary trees","Buy milk"],"lines":["Learn how to invert binary trees","Buy milk"]}',
        200,
        headers: {'content-type': 'application/json'},
      );
    });

    final api = HttpTaskApi(baseUrl: 'http://localhost:25664', client: client);
    await api.addTask('Buy milk');

    expect(requestBody, '{"task":"Buy milk"}');
  });

  test('deletes a task through the canonical endpoint', () async {
    late Uri requestedUri;
    final client = MockClient((request) async {
      requestedUri = request.url;
      return http.Response(
        '{"tasks":["Learn how to invert binary trees"],"lines":["Learn how to invert binary trees"]}',
        200,
        headers: {'content-type': 'application/json'},
      );
    });

    final api = HttpTaskApi(baseUrl: 'http://localhost:25664', client: client);
    await api.removeTask('Buy milk');

    expect(requestedUri.path, '/api/tasks');
    expect(requestedUri.queryParameters, {'task': 'Buy milk'});
  });
}
