import 'package:flutter/material.dart';

import '../code/task_list_controller.dart';
import '../contracts/task.dart';
import '../contracts/task_api.dart';

class TodoListPage extends StatefulWidget {
  final TaskApi api;

  const TodoListPage({super.key, required this.api});

  @override
  State<TodoListPage> createState() => _TodoListPageState();
}

class _TodoListPageState extends State<TodoListPage> {
  final TextEditingController _taskController = TextEditingController();
  TaskListViewModel _viewModel = const TaskListViewModel(tasks: []);

  @override
  void initState() {
    super.initState();
    _refreshTasks();
  }

  Future<void> _refreshTasks() async {
    final nextViewModel = await loadTasks(widget.api);
    if (!mounted) {
      return;
    }
    setState(() {
      _viewModel = nextViewModel;
    });
  }

  Future<void> _addTask() async {
    final nextViewModel = await addTask(_taskController.text, widget.api);
    _taskController.clear();
    if (!mounted) {
      return;
    }
    setState(() {
      _viewModel = nextViewModel;
    });
  }

  Future<void> _removeTask(Task task) async {
    final nextViewModel = await removeTask(task.id, widget.api);
    if (!mounted) {
      return;
    }
    setState(() {
      _viewModel = nextViewModel;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Todo List')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _taskController,
              decoration: const InputDecoration(labelText: 'Task'),
            ),
            const SizedBox(height: 12),
            ElevatedButton(onPressed: _addTask, child: const Text('Add task')),
            if (_viewModel.errorMessage != null) ...[
              const SizedBox(height: 12),
              Text(_viewModel.errorMessage!),
            ],
            const SizedBox(height: 12),
            Expanded(
              child: ListView(
                children: _viewModel.tasks
                    .map(
                      (task) => ListTile(
                        title: Text(task.text),
                        trailing: IconButton(
                          key: Key('remove-${task.id}'),
                          onPressed: () => _removeTask(task),
                          icon: const Icon(Icons.delete_outline),
                        ),
                      ),
                    )
                    .toList(),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
