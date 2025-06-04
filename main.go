package main

import (
	"encoding/json"
	"os"
	"strings"
)

// Start of the Aardvark application
func main() {
	choice := os.Args[1]
	environment := os.Args[2]
	fqdn = os.Args[3]
	var dot []string

	metadata := read("/data/automation/resources/" + environment + ".json")
	json.Unmarshal([]byte(metadata), &target)
	// createTable()

	if strings.Contains(fqdn, "/") {
		slash := strings.Split(fqdn, "/")
		slug = slash[1]
		hierarchy = slash[0]
	} else if strings.HasPrefix(fqdn, "www") {
		dot = strings.Split(fqdn, ".")
		slug = dot[1]
		hierarchy = dot[2]
	} else {
		dot = strings.Split(fqdn, ".")
		slug = dot[0]
		hierarchy = dot[1]
	}

	if !fileExists(target["workspace"] + "resources/" + target["sites"]) {
		document(target["workspace"]+"resources/"+target["sites"], []byte(getSites()))
	}

	siteID = getID(string(read(target["workspace"] + "resources/" + target["sites"])))

	if !fileExists(target["workspace"] + "temp/" + slug + ".csv") {
		document(target["workspace"]+"temp/"+slug+".csv", []byte(getPlugins()))
	}

	switch choice {
	case "-x":
		banner("The " + choice + " switch is being called to archive the " + fqdn + " blog site.")

		banner("Writing the archive event to the " + database + " database")
		// insertRow("archived")
	case "-a":
		source, destination = target["assets"]+siteID+"/", target["vault"]+siteID+"/"
		banner("Exporting the " + fqdn + " database")
		exportDB()

		banner("Exporting the " + fqdn + " users")
		exportUsers()

		banner("Exporting the " + fqdn + " assets")
		copyAssets()
		// direct(confirm(), "ac")

		err := zipFiles(slug+".zip", slug+".json", slug+".sql", slug+".csv", target["vault"]+siteID)
		inspect(err)

		banner("Moving " + slug + ".zip to the Jenkins workspace folder")
		execute("-e", "mv", slug+".zip", target["jenkins"])

		banner("Writing the archive event to the " + database + " database")
		// insertRow("archived")
	case "-r":
		source, destination = target["vault"]+siteID+"/", target["assets"]+siteID+"/"
		if err := unzip(slug+".zip", target["assets"]+siteID+"/"); err != nil {
			banner("Error unzipping file.")
			inspect(err)
			os.Exit(1)
		}
		banner("Successfully unzipped: " + slug + ".zip")

		banner("Importing the " + fqdn + " database")
		importDB()

		banner("Importing the " + fqdn + " assets (dry run)")
		copyAssetsDR()
		direct(confirm(), "ac")

		banner("Fixing HTTP References (dry run)")
		fixProtocolDR()
		direct(confirm(), "hf")

		banner("Writing the restore event to the " + database + " database")
		// insertRow("restored")
	case "-d":
		execute("-v", "wp", "site", "delete", siteID, "--path="+target["wordpress"], "--yes")
		banner("Writing the delete event to the " + database + " database")
		// insertRow("deleted")
	}

	banner("Flushing the WordPress cache")
	flushCache()
}
