package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	sourcePath := flag.String("p", "", "The path to source code of the application")
	outputPath := flag.String("o", "./data/go/", "The output path to save the results in")
	loglevel := flag.String("l", "default", "The logging level")
	printTree := flag.Bool("s", false, "Print the tree of the application")
	isMono := flag.Bool("m", false, "analyze a monolithic application")
	flag.Parse()

	appName := filepath.Base(*sourcePath)

	if *printTree {
		// Print tree for application
		outString, _, _, _, _ := Walk(*sourcePath, 0)
		logger.Info("Tree for ", appName, ":")
		logger.Info(outString)
		os.Exit(0)
	}

	// Set logger
	err := os.MkdirAll(filepath.Join(*outputPath, appName), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	logger.UpdateFile(filepath.Join(*outputPath, appName, "logs.log"))
	logger.SetLevel(*loglevel)

	// Find services
	logger.Info("Processing application: ", appName)
	serviceMap := make(map[string]string)
	if !*isMono {
		serviceFinder := NewServiceFinder(*sourcePath)
		serviceFinder.CreateRoot()
		_, services := serviceFinder.GetServices()

		logger.Info("Found ", len(services), " services")
		logger.Info("Services: ")
		for _, service := range services {
			logger.Info("Adding service: ", service)
			name := appName + "-" + filepath.Base(service)
			serviceMap[name] = service
		}
	} else {
		serviceMap[appName] = *sourcePath
	}
	logger.Debug("Found ", len(serviceMap), " services")

	// Process each service
	allObjects := []*Object_{}
	allExecutables := []Executable_{}
	for name, path := range serviceMap {
		logger.Debug("Working on service: ", name, " at ", path)
		microParser, err := ParseProject(path, name)
		if err != nil {
			logger.Error("Failed to parse service ", name, " : ", err)
			continue
		}
		allObjects = append(allObjects, Values(microParser.objects)...)
		allExecutables = append(allExecutables, Values(microParser.executables)...)
	}
	// Show results
	logger.Info("Detected ", len(allObjects), " structs/interfaces")
	logger.Info("Detected ", len(allExecutables), " methods/functions")
	logger.Info("Detected ", len(serviceMap), " microservices")

	// Save data
	// Implement the saveData function
	savePath := filepath.Join(*outputPath, appName)
	logger.Info("Saving data to ", savePath)
	saveData(allObjects, allExecutables, savePath)
}

func saveData(objects []*Object_, executables []Executable_, path string) {
	var filename, fileSavePath string
	var file []byte
	// Save objects
	filename = "typeData.json"
	fileSavePath = filepath.Join(path, filename)
	logger.Debug("Saving class data in " + fileSavePath)
	file, _ = json.MarshalIndent(objects, "", " ")
	_ = os.WriteFile(fileSavePath, file, 0644)
	// Save executables
	filename = "methodData.json"
	fileSavePath = filepath.Join(path, filename)
	logger.Debug("Saving class data in " + fileSavePath)
	file, _ = json.MarshalIndent(executables, "", " ")
	_ = os.WriteFile(fileSavePath, file, 0644)
}

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	// from https://cs.opensource.google/go/x/exp/+/9ff063c7:maps/maps.go;l=20
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
