package main

import "fmt"

type UserInput struct {
	versionedPackages []string
	systemPackages    []string
	ignoredPackages   []string
	dryRun            bool
}

func (userInput *UserInput) HandleSubscribedPackageChanges(hostConfig Host, config Config, baseDir string) {
	if hostConfig.SubscribeTo != nil {
		for _, subscribedTo := range hostConfig.SubscribeTo {
			subscribedToConfig, _ := GetHostConfig(subscribedTo, config)
			subscribedToPackages, _ := ReadPackagesFromFile(subscribedToConfig.File, baseDir)
			comparisionResult := ComparePackages(subscribedToPackages, userInput.systemPackages)
			filteredResult := FilterComparisonResult(userInput.ignoredPackages, comparisionResult)
			userInput.askForUserInput(filteredResult)
		}
	}
}

func (userInput *UserInput) HandleHostPackagesChange() {
	comparisionResult := ComparePackages(userInput.versionedPackages, userInput.systemPackages)
	fmt.Println("Comparing installed packages with versioned packages")
	userInput.askForUserInput(comparisionResult)
}

func (userInput *UserInput) askForUserInput(comparisionResult CompareResult) {
	fmt.Println("Packages to remove:")
	result := askForUserInput("remove", comparisionResult.Removed, userInput.dryRun)
	userInput.systemPackages = removeElementsFromSlice(userInput.systemPackages, result)

	fmt.Println("Packages to add")
	result = askForUserInput("add", comparisionResult.Added, userInput.dryRun)
	userInput.systemPackages = addElementsToSlice(userInput.systemPackages, result)
}

func addElementsToSlice(s []string, toAdd []string) []string {
	for _, v := range toAdd {
		if index := getElementIndex(s, v); index == -1 {
			s = append(s, v)
		}
	}
	return s
}

// returns accepted packages
func askForUserInput(text string, packages []string, dryRun bool) []string {
	var result []string
	for _, p := range packages {
		fmt.Println(fmt.Sprintf("%s: %s", text, p))
		fmt.Println("[y/N]")
		if askForConfirmation() {
			if !dryRun {
				// todo add logic to remove or add package
			}
			result = append(result, p)
		}
	}
	return result
}

func removeElementsFromSlice(s []string, toRemove []string) []string {
	for _, v := range toRemove {
		if index := getElementIndex(s, v); index > -1 {
			s = deleteElementByIndex(s, index)
		}
	}
	return s
}
