package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const SEP = "$$$"
const PKGSEP = "$$"
const STRUCTSEP = "."
const CONTENTSEP = "\n\n"

type MicroParser struct {
	serviceRoot     string
	serviceName     string
	structs         map[string]map[string]Object_
	objects         map[string]*Object_
	executables     map[string]Executable_
	objectMethodMap map[string][]string
}

func (microParser *MicroParser) UpdateStructs() {
	for objectName, methodList := range microParser.objectMethodMap {
		object := microParser.objects[objectName]
		fullContent := []string{object.Content}
		for _, executableName := range methodList {
			executable := microParser.executables[executableName]
			fullContent = append(fullContent, executable.Content)
		}
		if len(fullContent) > 1 {
			logger.Debug("Updating struct:", object.FullName)
			object.Content = strings.Join(fullContent, CONTENTSEP)
		}

	}
}

func (microParser *MicroParser) ParseStructs(files map[string]string) {
	for _, file := range files {
		logger.Debug("Parsing file:", file)
		microParser.ParseStructsInFile(file)
	}
}

func (microParser *MicroParser) ParseStructsInFile(filePath string) {
	fset := token.NewFileSet()

	// Parse the file
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	packageName := f.Name.Name
	pkgFullName := getPkgFullName(packageName, filePath, microParser.serviceRoot)
	logger.Debug("Package:", packageName)
	pkgMap, ok := microParser.structs[pkgFullName]
	if !ok {
		pkgMap = make(map[string]Object_)
		microParser.structs[pkgFullName] = pkgMap
	}

	// Analyze the AST and search for structs and interfaces
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			// Print struct name
			if _, ok := x.Type.(*ast.StructType); ok {
				object := Object_{
					IsInterface:  false,
					IsAnnotation: false,
					SimpleName:   x.Name.Name,
					FullName:     getStructFullName(*x, pkgFullName),
					FilePath:     filePath,
					ServiceName:  microParser.serviceName,
					Content:      getCodeSnippet(fset, x.Pos(), x.End()),
				}
				microParser.objectMethodMap[object.FullName] = []string{}
				logger.Debug("Struct:", object.FullName)
				microParser.objects[object.FullName] = &object
				pkgMap[object.SimpleName] = object
			} else if _, ok := x.Type.(*ast.InterfaceType); ok {
				object := Object_{
					IsInterface:  true,
					IsAnnotation: false,
					SimpleName:   x.Name.Name,
					FullName:     getStructFullName(*x, pkgFullName),
					FilePath:     filePath,
					ServiceName:  microParser.serviceName,
					Content:      getCodeSnippet(fset, x.Pos(), x.End()),
				}
				microParser.objectMethodMap[object.FullName] = []string{}
				logger.Debug("Interface:", object.FullName)
				microParser.objects[object.FullName] = &object
			}
			//logger.Debug("Struct:", x.Fields.List[0].Type.(*ast.Ident).Name) //, x.Fields.List[0].Names[0].Name)
		}
		return true
	})
}

func (microParser *MicroParser) ParseFunctions(files map[string]string) {
	for _, file := range files {
		logger.Debug("Parsing file for functions/methods:", file)
		microParser.ParseFunctionsInFile(file)
	}
}

func (microParser *MicroParser) AddMethod(x *ast.FuncDecl, object Object_, fset *token.FileSet) Executable_ {
	executable := Executable_{
		SimpleName:  x.Name.Name,
		FullName:    getMethodFullName(x, object.FullName),
		ServiceName: microParser.serviceName,
		Content:     getCodeSnippet(fset, x.Pos(), x.End()),
		ParentName:  object.FullName,
	}
	microParser.executables[executable.FullName] = executable
	microParser.objectMethodMap[object.FullName] = append(microParser.objectMethodMap[object.FullName], executable.FullName)
	return executable
}

func (microParser *MicroParser) AddFunction(x *ast.FuncDecl, pkgFullName string, fset *token.FileSet) Executable_ {
	executable := Executable_{
		SimpleName:  x.Name.Name,
		FullName:    getFuncFullName(x, pkgFullName),
		ServiceName: microParser.serviceName,
		Content:     getCodeSnippet(fset, x.Pos(), x.End()),
		ParentName:  "",
	}
	microParser.executables[executable.FullName] = executable
	return executable
}

func (microParser *MicroParser) ParseFunctionsInFile(filePath string) {
	fset := token.NewFileSet()

	// Parse the file
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	packageName := f.Name.Name
	pkgFullName := getPkgFullName(packageName, filePath, microParser.serviceRoot)
	logger.Debug("Package:", packageName)
	pkgMap, ok := microParser.structs[pkgFullName]
	if !ok {
		pkgMap = make(map[string]Object_)
		microParser.structs[pkgFullName] = pkgMap
	}

	// Analyze the AST and search for functions and methods
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Print function name
			if x.Recv != nil {
				// Method
				var executable Executable_
				methodName := x.Name.Name
				var structName string
				starExpr, ok := x.Recv.List[0].Type.(*ast.StarExpr)
				if ok {
					structName = starExpr.X.(*ast.Ident).Name
				} else {
					structName = x.Recv.List[0].Type.(*ast.Ident).Name
				}
				object, ok := pkgMap[structName]
				if !ok {
					logger.Warn("Struct ", structName, " not found for method ", methodName)
					executable = microParser.AddFunction(x, pkgFullName, fset)
					logger.Debug("Method as function:", executable.FullName)
				} else {
					executable = microParser.AddMethod(x, object, fset)
					logger.Debug("Method:", executable.FullName)
				}
			} else {
				// Function
				executable := microParser.AddFunction(x, pkgFullName, fset)
				logger.Debug("Function:", executable.FullName)
				//logger.Debug("Function:", x.Name.Name, " with fullName ", getFuncFullName(x, pkgFullName))
			}
		}
		return true
	})
}

func findGoFiles(root string) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			files[info.Name()] = path
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func ParseProject(root string, serviceName string) (*MicroParser, error) {
	files, err := findGoFiles(root)
	if err != nil {
		return nil, err
	}
	microParser := &MicroParser{
		serviceRoot:     root,
		serviceName:     serviceName,
		structs:         make(map[string]map[string]Object_),
		objects:         make(map[string]*Object_),
		executables:     make(map[string]Executable_),
		objectMethodMap: make(map[string][]string),
	}
	microParser.ParseStructs(files)
	microParser.ParseFunctions(files)
	microParser.UpdateStructs()

	return microParser, nil
}

func getFuncFullName(fn *ast.FuncDecl, pkgName string) string {
	// Get the function signature
	return pkgName + SEP + fn.Name.Name + "()"
}

func getStructFullName(st ast.TypeSpec, pkgName string) string {
	return pkgName + SEP + st.Name.Name
}

func getMethodFullName(fn *ast.FuncDecl, structName string) string {
	return structName + STRUCTSEP + fn.Name.Name + "()"
}

func getCodeSnippet(fset *token.FileSet, start, end token.Pos) string {
	fileContent, err := os.ReadFile(fset.File(start).Name())
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	return string(fileContent[start-1 : end])
}

func getPkgFullName(packageName, filePath, projectRoot string) string {
	prefix := filepath.Dir(filePath)
	pkgFullName := strings.Replace(strings.Replace(prefix, projectRoot, "", 1), string(os.PathSeparator), "/", -1) + PKGSEP + packageName
	if strings.HasPrefix(pkgFullName, "/") {
		pkgFullName = pkgFullName[1:]
	}
	//if strings.HasSuffix(pkgFullName, "/") {
	//	pkgFullName = pkgFullName[:len(pkgFullName)-1]
	//}
	return pkgFullName
}

func FindGoDirectories(root string) ([]string, error) {
	dirs := make(map[string]struct{})

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			dir := filepath.Dir(path)
			dirs[dir] = struct{}{}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	uniqueDirs := make([]string, 0, len(dirs))
	for dir := range dirs {
		uniqueDirs = append(uniqueDirs, dir)
	}

	return uniqueDirs, nil
}
