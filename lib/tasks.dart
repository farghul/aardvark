import 'package:aardvark/tools.dart';

dynamic siteNumber, slug, unmarshal, url;

// Query WordPress for a list of all sites in the targeted environment
String getSiteList() {
  var query = execute('-r', 'wp', [
    'site',
    'list',
    '--fields=blog_id,url',
    '--ssh=${unmarshal['user']}@${unmarshal['server']}:${unmarshal['install']}',
    '--url=${unmarshal['address']}',
    '--skip-plugins',
    '--skip-packages',
    '--format=csv',
  ]);

  String result = query.stdout;
  result = result.replaceFirst('blog_id,url\n', '');
  result = result.replaceAll('https://', '');
  result = result.replaceAll('http://', '');
  result = result.replaceAll('/\n', ',');
  result = result.replaceAll('http://', '');
  return query.stdout;
}

// Search the blog list to find the SiteNumber that matches the supplied URL
String getSiteNumber(String environment, String list) {
  int counter = 0;
  String result = '';
  List<String> blogs = list.split(',');
  if (environment == 'production' || environment == 'test') {
    for (var element in blogs) {
      counter++;
      if (element.contains(slug)) {
        result = blogs[counter - 1];
      }
    }
  } else {
    for (var element in blogs) {
      if (element == '${unmarshal['address']}$slug') {
        result = blogs[counter - 1];
      }
    }
  }
  return result;
}

// Query WordPress for a list of plugins installed relative to a specific site, and their current version
String getPluginList() {
  var result = execute('-r', 'wp', [
    'plugin list --status=active --fields=name,version --ssh=${unmarshal['user']}@${unmarshal['server']}:${unmarshal['install']} --url=${unmarshal['address']} --skip-plugins --skip-themes --skip-packages --format=csv',
  ]);
  return result.stdout;
}

// Query WordPress for a list of themes installed relative to a specific site, and their current version
String getThemeList() {
  var result = execute('-r', 'wp', [
    'theme list --status=active --fields=name,version --ssh=${unmarshal['user']}@${unmarshal['server']}:${unmarshal['install']} --url=${unmarshal['address']} --skip-plugins --skip-themes --skip-packages --format=csv',
  ]);
  return result.stdout;
}

// Export the source WordPress site database to an SQL file
void exportDatabase() {
  var inner = execute('-d', 'wp', [
    'db tables --all-tables-with-prefix --ssh=${unmarshal['user']}@${unmarshal['server']}:${unmarshal['install']} --url=${unmarshal['address']} --skip-plugins --skip-themes --skip-packages --format=csv',
  ]);
  execute('-d', 'wp', [
    'db export --tables=${inner.stdout} ${unmarshal['ephemeral']}tables.sql --ssh=${unmarshal['user']}@${unmarshal['server']}:${unmarshal['install']}',
  ]);
}

// Create a user export file in CSV format
String exportUserList() {
  var result = execute('-r', 'wp', [
    'user list --ssh=${unmarshal['user']}@${unmarshal['server']}:${unmarshal['install']} --url=${unmarshal['address']} --skip-plugins --skip-themes --skip-packages --format=csv',
  ]);
  return result.stdout;
}

// Copy WordPress site assets to a new location
void copyAssets(String source, String destination) {
  execute('-d', 'cd', [destination]);
  execute('-d', 'scp', ['-r $source $destination']);
  execute('-d', 'tar', ['-cvf assets.tar $destination$siteNumber']);
}

// Flush the WordPress and Composer caches
void flushCaches() {
  execute('-d', 'wp', [
    'wp cache flush --ssh=${unmarshal['user']}@${unmarshal['server']}:${unmarshal['install']}',
  ]);
  execute('-d', 'composer', ['clear-cache']);
}
