package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func GetInstalledPackages(hostPackageFilePath *string) ([]string, error) {
	var packages []string
	if hostPackageFilePath != nil {
		packages, _ = ReadPackagesFromFile(*hostPackageFilePath, "")
	} else {
		cmd := exec.Command("yaourt", "-Qqe")
		result, err := cmd.Output()
		if err != nil {
			return []string{}, nil
		}
		resultAsString := string(result)
		packages = strings.Split(resultAsString, "\n")
	}
	return ClearPackages(packages), nil
}

// remove empty strings from the packages
func ClearPackages(packages []string) []string {
	var clearedPackages []string
	for _, v := range packages {
		if len(v) > 1 {
			clearedPackages = append(clearedPackages, v)
		}
	}
	return clearedPackages
}

func CommitPackageChanges(userInput UserInput, hostConfig Host, baseDir string, dryRun bool) {
	if !dryRun {
		filePath := baseDir + "/" + hostConfig.File
		systemPackages := []byte(strings.Join(userInput.systemPackages, "\n"))
		err := ioutil.WriteFile(filePath, systemPackages, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to save versioned packages:", err)
		}
		ignoredPackages := []byte(strings.Join(userInput.ignoredPackages, "\n"))
		filePath = baseDir + "/" + hostConfig.IgnoreFile
		err = ioutil.WriteFile(filePath, ignoredPackages, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to save ignored packages:", err)
		}
	} else {
		fmt.Println("New versioned packages:")
		for _, v := range userInput.systemPackages {
			fmt.Println(v)
		}
		fmt.Println("New ignored packages:")
		for _, v := range userInput.ignoredPackages {
			fmt.Println(v)
		}
	}
}
