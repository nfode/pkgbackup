package main

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"os"
	"os/exec"
	"strings"
)

func openEditor(path string) error {
	editor := os.Getenv("EDITOR")
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func WriteTmpFile(comparisionResult CompareResult) string {
	var result []string
	if len(comparisionResult.Added) > 0 {
		result = append(result, "# Added:")
		for _, v := range comparisionResult.Added {
			result = append(result, "install "+v)
		}
	}
	if len(comparisionResult.Removed) > 0 {
		result = append(result, "# Remove:")
		for _, v := range comparisionResult.Removed {
			result = append(result, "remove "+v)
		}
	}
	result = append(result, footer())
	resultAsString := strings.Join(result, "\n")
	id, _ := uuid.NewV4()
	filePath := "/tmp/pkbackup_" + id.String() + ".txt"
	fmt.Println(resultAsString)
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		fmt.Println("creating tmp file failed: " + err.Error())
	}
	f.WriteString(resultAsString)
	f.Sync()
	return filePath
}

func AskUser(comparisionResult CompareResult) error {
	path := WriteTmpFile(comparisionResult)
	defer os.Remove(path)
	return openEditor(path)
}

func footer() string {
	return `
# Commands:
# 	i, install <package> = install package
# 	r, remove <package> = remove package
# 	a, add <package> = add package to package backup file
		`
}
