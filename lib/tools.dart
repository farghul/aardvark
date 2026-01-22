import 'dart:io';
import 'package:process_run/shell.dart';

var shell = Shell();

// Run standard terminal commands
dynamic execute(String variation, String command, List<String> arguments) {
  switch (variation) {
    // Direct execution
    case '-d':
      shell.runExecutableArguments(command, arguments);
      break;
    // Return a value
    case '-r':
      var result = shell.runExecutableArgumentsSync(command, arguments);
      return result.stdout;
    default:
  }
}

// Provide an informational message
void printBanner(String message) {
  print('*** $message ***');
}

// Read any file and return the contents as a String variable
String readFile(File document) {
  String contents = document.readAsStringSync();
  return contents;
}

// Write a passed variable to a named file
void writeFile(File document, String content) {
  document.writeAsStringSync(content);
}

// Delete a file
void deleteFile(File document) {
  document.deleteSync();
}

// Delete the contents a folder
void emptyDirectory(String destination) async {
  var target = Directory(destination);
  await for (var entity in target.list(recursive: true, followLinks: false)) {
    execute('-d', 'rm', [entity.path]);
  }
}
