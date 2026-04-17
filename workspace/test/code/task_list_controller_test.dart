import 'package:mocktail/mocktail.dart';
import 'package:todo_list/code/task_list_controller.dart';
import 'package:todo_list/contracts/task.dart';
import 'package:todo_list/contracts/task_api.dart';
import 'package:todo_list/contracts/task_list_response.dart';
import 'package:test/test.dart';

class MockTaskApi extends Mock implements TaskApi {}

void main() {
  test('loads the current task list', () async {
    final api = MockTaskApi();
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: [
          Task(
            id: '11111111-1111-1111-1111-111111111111',
            text: 'Learn how to invert binary trees',
          ),
          Task(id: '22222222-2222-2222-2222-222222222222', text: 'Buy milk'),
        ],
      ),
    );

    final result = await loadTasks(api);

    expect(result.tasks.map((task) => task.text).toList(), [
      'Learn how to invert binary trees',
      'Buy milk',
    ]);
    expect(result.errorMessage, isNull);
  });

  test('trims submitted tasks before calling the REST API', () async {
    final api = MockTaskApi();
    when(() => api.createTask('Buy milk')).thenAnswer(
      (_) async => const Task(
        id: '22222222-2222-2222-2222-222222222222',
        text: 'Buy milk',
      ),
    );
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: [
          Task(
            id: '11111111-1111-1111-1111-111111111111',
            text: 'Learn how to invert binary trees',
          ),
          Task(id: '22222222-2222-2222-2222-222222222222', text: 'Buy milk'),
        ],
      ),
    );

    final result = await addTask('  Buy milk  ', api);

    expect(result.tasks.map((task) => task.text).toList(), [
      'Learn how to invert binary trees',
      'Buy milk',
    ]);
    verify(() => api.createTask('Buy milk')).called(1);
  });
}
