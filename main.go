package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
)

var (
	app             = kingpin.New("pkgbackup", "generic help")
	configFile      = app.Flag("configFile", "path to the configFile file").File()
	hostPackageFile = app.Flag("hostPackageFile", "alternative host packages").File()
	dryRunFlag      = app.Flag("dryRun", "run the tool without installing anything").Bool()
	syncCommand     = app.Command("sync", "export package list")
	hostname        = os.Getenv("HOST")
)

func main() {
	res := kingpin.MustParse(app.Parse(os.Args[1:]))
	config := ParseConfigFile(*configFile)
	switch res {
	case syncCommand.FullCommand():
		sync(config, path.Dir((*configFile).Name()))
	}
}
func sync(config Config, baseDir string) {
	hostConfig, err := GetHostConfig(hostname, config)
	versionedPackages, err := ReadPackagesFromFile(hostConfig.File, baseDir)
	systemPackages, err := GetInstalledPackages()
	if err != nil {
		fmt.Println("reading existing packages failed: " + err.Error())
	}
	ignoredPackages, err := ReadPackagesFromFile(hostConfig.IgnoreFile, baseDir)
	if err != nil {
		fmt.Println(fmt.Printf("Failed to read ignore file for host %v: %v", hostname, err.Error()))
	}
	userInput := UserInput{dryRun: *dryRunFlag, versionedPackages: versionedPackages, systemPackages: systemPackages, ignoredPackages: ignoredPackages}
	userInput.HandleHostPackagesChange()
	userInput.HandleSubscribedPackageChanges(hostConfig, config, baseDir)
	CommitPackageChanges(userInput, hostConfig, baseDir)
}






