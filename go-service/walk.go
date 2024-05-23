package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Walk(sourcePath string, prefix int) (string, bool, bool, bool, bool) {
	//var output strings.Builder
	var hasGoFile, hasModFile, hasMainFile, hasMain bool

	files, err := os.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	folderName := filepath.Base(sourcePath)
	var currentString string
	var subStrings []string
	indent := strings.Repeat("|   ", prefix)
	nextIndent := strings.Repeat("|   ", prefix+1)
	currentString = indent + "|---" + folderName + "/"
	for _, file := range files {
		if file.IsDir() {
			//output.WriteString(indent + "|---" + info.Name() + "/\n")
			subOutput, subHasGoFile, subHasModFile, subHasMainFile, subHasMain := Walk(filepath.Join(sourcePath, file.Name()), prefix+1)
			//output.WriteString(subOutput)
			if subOutput != "" {
				subStrings = append(subStrings, subOutput)
			}
			hasGoFile = hasGoFile || subHasGoFile
			hasModFile = hasModFile || subHasModFile
			hasMainFile = hasMainFile || subHasMainFile
			hasMain = hasMain || subHasMain
		} else {
			if file.Name() == "go.mod" {
				hasModFile = true
				fileString := nextIndent + "|---" + file.Name()
				subStrings = append(subStrings, fileString)
			}
			if file.Name() == "main.go" {
				hasMainFile = true
			}
			if strings.HasSuffix(file.Name(), ".go") {
				if HasMainFunction(filepath.Join(sourcePath, file.Name())) {
					hasMain = true
				}
				fileString := nextIndent + "|---" + file.Name()
				subStrings = append(subStrings, fileString)
				hasGoFile = true
			}
		}
	}
	if hasMain || hasMainFile || hasModFile {
		currentString += " $"
	}
	if hasGoFile {
		currentString += " *"
	}
	if !(hasMain || hasMainFile || hasModFile || hasGoFile) {
		return "", false, false, false, false
	}
	output := currentString
	if len(subStrings) > 0 {
		output += "\n" + strings.Join(subStrings, "\n")
	}
	return output, hasGoFile, hasModFile, hasMainFile, hasMain
}
