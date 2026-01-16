import 'dart:convert';
import 'dart:io';

import 'package:args/args.dart';
import 'package:aardvark/tools.dart';
import 'package:aardvark/tasks.dart';

// Start of the archiving application
void main(List<String> arguments) {
  var parser = ArgParser();
  parser.addOption(
    'environment',
    abbr: 'e',
    allowed: ['development', 'production', 'staging', 'test'],
    mandatory: true,
  );

  ArgResults argResults = parser.parse(arguments);
  final paths = argResults.rest;

  String location = '/data/automation/checkouts/dac/jsons/';
  File metadata = File('$location${argResults.option('environment')}.json');
  unmarshal = jsonDecode(readFile(metadata));

  String url = paths[0];
  url = url.replaceAll('https://', '');
  url = url.replaceAll('http://', '');
  List<String> halves, parts;

  if (url.contains('/')) {
    halves = url.split('/');
    parts = url.split('.');
    if (halves[1].isEmpty) {
      slug = parts[0];
    } else {
      slug = halves[1];
    }
  } else {
    parts = url.split('.');
    slug = parts[0];
  }

  writeFile(File('${unmarshal['lists']}${unmarshal['sites']}'), getSites());

  siteID = getID(
    '${argResults.option('environment')}',
    readFile(File('${unmarshal['lists']}${unmarshal['sites']}')),
  );

  emptyDir('${unmarshal['ephemeral']}');

  writeFile(File('${unmarshal['ephemeral']}id.txt'), siteID);
  writeFile(File('${unmarshal['ephemeral']}plugins.csv'), getPlugins());
  writeFile(File('${unmarshal['ephemeral']}themes.csv'), getThemes());

  execute('-d', 'cp', [
    '/data/www-app/${unmarshal['title']}/current/composer.lock',
    '${unmarshal['ephemeral']}',
  ]);

  banner('Exporting the $slug database');
  exportDatabase();

  writeFile(File('${unmarshal['ephemeral']}users.csv'), exportUsers());
  banner('Exporting the $slug users');

  banner('Exporting the $slug assets');
  copyAssets('${unmarshal['assets']}', '${unmarshal['ephemeral']}');

  banner('Flushing the WordPress cache');
  flushCaches();
}
