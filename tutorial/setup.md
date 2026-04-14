# Setup

Keep the repository root for shared files like `README.md`, `LICENSE`, `.gitignore`, `.github/`, `justfile`, and `tutorial/`.

Put all Flutter code inside a single `workspace/` folder.

From the repository root, run each setup command and checkpoint it before moving to the next one:

```bash
flutter create --platforms=web,android,ios,macos,windows,linux --org com.intrepion --project-name todo_list workspace
git add --all
git commit --message "flutter create --platforms=web,android,ios,macos,windows,linux --org com.intrepion --project-name todo_list workspace"

rm workspace/test/widget_test.dart
git add --all
git commit --message "rm workspace/test/widget_test.dart"

mkdir -p workspace/integration_test
git add --all
git commit --message "mkdir -p workspace/integration_test"

mkdir -p workspace/lib/contracts
git add --all
git commit --message "mkdir -p workspace/lib/contracts"

mkdir -p workspace/lib/code
git add --all
git commit --message "mkdir -p workspace/lib/code"

mkdir -p workspace/lib/adapter
git add --all
git commit --message "mkdir -p workspace/lib/adapter"

mkdir -p workspace/test/code
git add --all
git commit --message "mkdir -p workspace/test/code"

mkdir -p workspace/test/adapter
git add --all
git commit --message "mkdir -p workspace/test/adapter"

(cd workspace && flutter pub add http)
just format
git add --all
git commit --message "(cd workspace && flutter pub add http)"

(cd workspace && flutter pub add --dev test)
just format
git add --all
git commit --message "(cd workspace && flutter pub add --dev test)"

(cd workspace && flutter pub add --dev mocktail)
just format
git add --all
git commit --message "(cd workspace && flutter pub add --dev mocktail)"

(cd workspace && flutter pub add --dev integration_test --sdk flutter)
just format
git add --all
git commit --message "(cd workspace && flutter pub add --dev integration_test --sdk flutter)"
```

When the full workspace is finished, it should contain these files:

```text
workspace/
  pubspec.yaml
  lib/
    contracts/
      task_api.dart
      task_list_response.dart
    code/
      task_list_controller.dart
    adapter/
      http_task_api.dart
      todo_list_page.dart
  test/
    code/
      task_list_controller_test.dart
    adapter/
      http_task_api_test.dart
      todo_list_page_test.dart
  integration_test/
    app_test.dart
  lib/main.dart
```

Before you try any run command, make sure Flutter can see a supported target:

```bash
just devices
```

For web, use the default web command:

```bash
just run
```

or, explicitly:

```bash
just run-web
```

On macOS for iOS, install CocoaPods first if you have not already:

```bash
sudo gem install cocoapods
```

Then open the simulator, list devices, and run the iOS app:

```bash
open -a Simulator
just devices
just run-ios
```

If you want to target a specific iOS simulator id or name, use:

```bash
just run-ios device="<ios-device-id-or-name>"
```

For Android, list available emulators, launch one, list devices again, and then run the Android app:

```bash
just emulators
flutter emulators --launch <emulator-id>
just devices
just --set api_base_url http://10.0.2.2:25664 run-android device="<android-device-id-or-name>"
```

For macOS desktop, use:

```bash
just run-macos
```

For Windows or Linux, run the matching command on that host platform:

```bash
just run-windows
just run-linux
```

After your first successful iOS run, if CocoaPods added shared iOS project files like these:

- `workspace/ios/Runner.xcodeproj/project.pbxproj`
- `workspace/ios/Runner.xcworkspace/contents.xcworkspacedata`
- `workspace/ios/Podfile.lock`

then run:

```bash
git add --all
git commit --message "Add iOS CocoaPods workspace files"
```

Do not commit local machine output like these:

- `workspace/ios/Pods/`
- `workspace/build/`
- `workspace/.dart_tool/`

For Android, a normal first run usually should not add shared tracked files. If it does change shared files under `workspace/android/`, review them carefully and commit only the project-level changes. Do not commit machine-specific files like:

- `workspace/android/local.properties`
- `workspace/.gradle/`
- `workspace/build/`
