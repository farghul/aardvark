package main

import (
	"os"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Constant declarations
const (
	automatic string = "\033[0m"
	yellowFG  string = "\033[33m"
	redBG     string = "\033[41m"
	database  string = "aardvark"
	table     string = "operations"
)

// Variables and map declarations
var (
	target map[string]string
	// db, _  = sql.Open("sqlite3", "file:"+target["workspace"]+"resources/"+database+".db?_txlock=immediate")
	// String variables
	fqdn, hierarchy, slug, source, siteID, destination string
)

// Create a new SQLite table if it doesn't exist
// func createTable() {
// 	db.Exec("CREATE TABLE IF NOT EXISTS " + table + " (id INTEGER PRIMARY KEY AUTOINCREMENT, slug TEXT NOT NULL, hierarchy TEXT NOT NULL, blog TEXT NOT NULL, operation TEXT NOT NULL, server TEXT NOT NULL, created REAL NOT NULL)")
// }

// Inser a row of data into the SQLite table
// func insertRow(operation string) {
// 	stmt, err := db.Prepare("INSERT INTO  " + table + "(slug, hierarchy, blog, operation, server, created) VALUES(?, ?, ?, ?, ?, ?)")
// 	inspect(err)
// 	defer stmt.Close()

// 	_, err = stmt.Exec(slug, hierarchy, siteID, operation, target["server"], time.Now())
// 	inspect(err)
// }

// Search the blog list to find the ID that matches the supplied URL
func getID(list string) string {
	var result string
	blogs := strings.Split(list, ",")
	for order, item := range blogs {
		if item == fqdn {
			result = blogs[order-1]
		}
	}

	return result
}

// Query WordPress for a list of all sites in the production environment
func getSites() string {
	query := execute("-c", "wp", "site", "list", "--fields=blog_id,url", "--path="+target["wordpress"], "--url="+target["address"], "--skip-plugins", "--skip-themes", "--skip-packages", "--format=csv")
	result := strings.Replace(string(query), "blog_id,url\n", "", 1)
	result = strings.ReplaceAll(result, "https://", "")
	result = strings.ReplaceAll(result, "http://", "")
	result = strings.ReplaceAll(result, "/\n", ",")
	result = strings.TrimSuffix(result, ",")

	return result
}

// Query WordPress for a list of plugins installed reletive to a specific site, and their current version
func getPlugins() string {
	query := execute("-c", "wp", "plugin", "list", "--status=active", "--fields=name,version", "--path="+target["wordpress"], "--url="+fqdn, "--skip-plugins", "--skip-themes", "--skip-packages", "--format=csv")
	result := strings.ReplaceAll(string(query), "/\n", ",")
	result = strings.TrimSuffix(result, ",")

	return result
}

// Export the source WordPress site database to an sql file
func exportDB() {
	inside := execute("-c", "wp", "db", "tables", "--all-tables-with-prefix", "--url="+fqdn, "--path="+target["wordpress"], "--skip-plugins", "--skip-themes", "--skip-packages", "--format=csv")
	result := strings.ReplaceAll(string(inside), "\n", ",")
	result = strings.TrimSuffix(result, ",")
	execute("-v", "wp", "db", "export", "--tables="+result, target["workspace"]+"temp/"+slug+".sql", "--path="+target["wordpress"])
}

// Create a user export file in JSON format
func exportUsers() {
	people := execute("-c", "wp", "user", "list", "--url="+fqdn, "--path="+target["wordpress"], "--skip-plugins", "--skip-themes", "--skip-packages", "--format=csv")
	inspect(os.WriteFile(target["workspace"]+"temp/"+slug+".json", people, 0666))
}

// Copy WordPress site assets to a new location
func copyAssets() {
	execute("-v", "rsync", "-a", source, destination)
	// fmt.Println("-v rsync -a " + source + " " + destination)
}

// Import WordPress SQL database tables
func importDB() {
	execute("-v", "wp", "db", "import", "/data/temp/"+slug+".sql", "--path="+target["wordpress"])
}

// Catch any lingering http addresses and convert them to https
func fixProtocol() {
	execute("-v", "wp", "search-replace", "--url="+fqdn, "--path="+target["wordpress"], "--all-tables-with-prefix", "http://", "https://")
}

// Flush the WordPress cache
func flushCache() {
	execute("-v", "wp", "cache", "flush", "--path="+target["wordpress"])
}

/* ----- Dry Run Functions ----- */

// Copy the site assets over using --dry-run
func copyAssetsDR() {
	execute("-v", "rsync", "-a", source, destination, "--stats", "--dry-run")
	// fmt.Println("-v rsync -a " + source + " " + destination + " --stats --dry-run")
}

// Catch any lingering http addresses using --dry-run
func fixProtocolDR() {
	execute("-v", "wp", "search-replace", "--url="+fqdn, "--path="+target["wordpress"], "--all-tables-with-prefix", "http://", "https://", "--dry-run")
}
