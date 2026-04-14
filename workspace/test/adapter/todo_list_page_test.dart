import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:todo_list/adapter/todo_list_page.dart';
import 'package:todo_list/contracts/task_api.dart';
import 'package:todo_list/contracts/task_list_response.dart';

class MockTaskApi extends Mock implements TaskApi {}

void main() {
  testWidgets('loads, adds, and removes tasks', (tester) async {
    final api = MockTaskApi();
    when(() => api.getTasks()).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees'],
        lines: ['Learn how to invert binary trees'],
      ),
    );
    when(() => api.addTask('Buy milk')).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees', 'Buy milk'],
        lines: ['Learn how to invert binary trees', 'Buy milk'],
      ),
    );
    when(() => api.removeTask('Buy milk')).thenAnswer(
      (_) async => const TaskListResponse(
        tasks: ['Learn how to invert binary trees'],
        lines: ['Learn how to invert binary trees'],
      ),
    );

    await tester.pumpWidget(MaterialApp(home: TodoListPage(api: api)));
    await tester.pumpAndSettle();

    expect(find.text('Learn how to invert binary trees'), findsOneWidget);

    await tester.enterText(find.byType(TextField), 'Buy milk');
    await tester.tap(find.text('Add task'));
    await tester.pumpAndSettle();

    expect(find.text('Buy milk'), findsOneWidget);

    await tester.tap(find.byKey(const Key('remove-Buy milk')));
    await tester.pumpAndSettle();

    expect(find.text('Buy milk'), findsNothing);
  });
}
