import 'package:mocktail/mocktail.dart';
import 'package:todo_list/code/task_list_controller.dart';
import 'package:todo_list/contracts/task_api.dart';
import 'package:todo_list/contracts/task_list_response.dart';
import 'package:test/test.dart';

class MockTaskApi extends Mock implements TaskApi {}

void main() {
  test('loads the current task list', () async {
    final api = MockTaskApi();
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees', 'Buy milk'],
        lines: ['Learn how to invert binary trees', 'Buy milk'],
      ),
    );

    final result = await loadTasks(api);

    expect(result.tasks, ['Learn how to invert binary trees', 'Buy milk']);
    expect(result.lines, ['Learn how to invert binary trees', 'Buy milk']);
    expect(result.errorMessage, isNull);
  });

  test('trims submitted tasks before calling the API', () async {
    final api = MockTaskApi();
    when(() => api.addTask('Buy milk')).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees', 'Buy milk'],
        lines: ['Learn how to invert binary trees', 'Buy milk'],
      ),
    );

    final result = await addTask('  Buy milk  ', api);

    expect(result.tasks, ['Learn how to invert binary trees', 'Buy milk']);
    verify(() => api.addTask('Buy milk')).called(1);
  });
}
