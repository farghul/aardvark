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

	metadata := read("/data/automation/resources/" + environment + ".json")
	json.Unmarshal([]byte(metadata), &target)
	slash := strings.Split(fqdn, "/")
	slug = slash[1]

	if !fileExists(target["workspace"] + "resources/" + target["sites"]) {
		document(target["workspace"]+"resources/"+target["sites"], []byte(getSites()))
	}

	siteID = getID(string(read(target["workspace"] + "resources/" + target["sites"])))

	if !fileExists(target["workspace"] + "temp/" + slug + ".csv") {
		document(target["workspace"]+"temp/"+slug+".csv", []byte(getPlugins()))
	}

	switch choice {
	case "-a":
		source, destination = target["assets"]+siteID+"/", target["vault"]+siteID+"/"
		execute("-e", "mkdir", destination)

		banner("Exported the " + fqdn + " database")
		exportDB()

		banner("Exported the " + fqdn + " users")
		exportUsers()

		banner("Exported the " + fqdn + " assets")
		copyAssets()

		err := zipFiles(slug+".zip", slug+".json", slug+".sql", slug+".csv", target["vault"]+siteID)
		inspect(err)
	case "-r":
		source, destination = target["vault"]+siteID+"/", target["assets"]+siteID+"/"
		if err := unzip(slug+".zip", target["assets"]+siteID+"/"); err != nil {
			banner("Error unzipping file.")
			inspect(err)
			os.Exit(1)
		}
		banner("Successfully unzipped: " + slug + ".zip")

		banner("Imported the " + fqdn + " database")
		importDB()

		banner("Imported the " + fqdn + " assets")
		copyAssets()

		banner("Fixed HTTP References")
		fixProtocol()
	case "-d":
		execute("-v", "wp", "site", "delete", siteID, "--path="+target["wordpress"], "--yes")
		banner("Writing the delete event to the " + database + " database")
	}

	banner("Flushing the WordPress cache")
	flushCache()
}
