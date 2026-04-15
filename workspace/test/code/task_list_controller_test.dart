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
          Task(id: 1, text: 'Learn how to invert binary trees'),
          Task(id: 2, text: 'Buy milk'),
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
    when(
      () => api.createTask('Buy milk'),
    ).thenAnswer((_) async => const Task(id: 2, text: 'Buy milk'));
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: [
          Task(id: 1, text: 'Learn how to invert binary trees'),
          Task(id: 2, text: 'Buy milk'),
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

  test('rejects blank submitted tasks without calling the API', () async {
    final api = MockTaskApi();

    final result = await addTask('   ', api);

    expect(result.errorMessage, 'Task must not be blank.');
    verifyNever(() => api.createTask(any()));
  });

  test('returns a friendly message when the task API is unavailable', () async {
    final api = MockTaskApi();
    when(() => api.getTasks()).thenThrow(Exception('boom'));

    final result = await loadTasks(api);

    expect(
      result.errorMessage,
      'Sorry, the task API is unavailable right now.',
    );
  });

  test('removes the chosen task through the REST API', () async {
    final api = MockTaskApi();
    when(() => api.deleteTask(2)).thenAnswer((_) async {});
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: [Task(id: 1, text: 'Learn how to invert binary trees')],
      ),
    );

    final result = await removeTask(2, api);

    expect(result.tasks.map((task) => task.text).toList(), [
      'Learn how to invert binary trees',
    ]);
    verify(() => api.deleteTask(2)).called(1);
  });
}
