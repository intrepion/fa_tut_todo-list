import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:todo_list/adapter/todo_list_page.dart';
import 'package:todo_list/contracts/task.dart';
import 'package:todo_list/contracts/task_api.dart';
import 'package:todo_list/contracts/task_list_response.dart';

class MockTaskApi extends Mock implements TaskApi {}

void main() {
  testWidgets('loads, adds, and removes tasks', (tester) async {
    final api = MockTaskApi();
    var loadCount = 0;
    when(() => api.getTasks()).thenAnswer((_) async {
      loadCount += 1;
      switch (loadCount) {
        case 1:
          return const TaskListResponse(
            tasks: [
              Task(
                id: '11111111-1111-1111-1111-111111111111',
                text: 'Learn how to invert binary trees',
              ),
            ],
          );
        case 2:
          return const TaskListResponse(
            tasks: [
              Task(
                id: '11111111-1111-1111-1111-111111111111',
                text: 'Learn how to invert binary trees',
              ),
              Task(
                id: '22222222-2222-2222-2222-222222222222',
                text: 'Buy milk',
              ),
            ],
          );
        default:
          return const TaskListResponse(
            tasks: [
              Task(
                id: '11111111-1111-1111-1111-111111111111',
                text: 'Learn how to invert binary trees',
              ),
            ],
          );
      }
    });
    when(() => api.createTask('Buy milk')).thenAnswer(
      (_) async => const Task(
        id: '22222222-2222-2222-2222-222222222222',
        text: 'Buy milk',
      ),
    );
    when(
      () => api.deleteTask('22222222-2222-2222-2222-222222222222'),
    ).thenAnswer((_) async {});

    await tester.pumpWidget(MaterialApp(home: TodoListPage(api: api)));
    await tester.pumpAndSettle();

    expect(find.text('Learn how to invert binary trees'), findsOneWidget);

    await tester.enterText(find.byType(TextField), 'Buy milk');
    await tester.tap(find.text('Add task'));
    await tester.pumpAndSettle();

    expect(find.text('Buy milk'), findsOneWidget);

    await tester.tap(
      find.byKey(const Key('remove-22222222-2222-2222-2222-222222222222')),
    );
    await tester.pumpAndSettle();

    expect(find.text('Buy milk'), findsNothing);
  });
}
