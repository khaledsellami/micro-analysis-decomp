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

type ServiceNode struct {
	path         string
	name         string
	parent       *ServiceNode
	children     []ServiceNode
	hasGo        bool
	hasMod       bool
	hasMainGo    bool
	hasMain      bool
	isService    bool
	CanBeService bool
	NGo          int
}

func (sn *ServiceNode) Update() {
	for _, child := range sn.children {
		if child.name == "go.mod" {
			sn.hasMod = true
		} else {
			if child.hasGo {
				sn.hasGo = true
			}
		}
		sn.CanBeService = sn.CanBeService || child.CanBeService
	}
	// TODO: Check if the node is a service
	sn.CanBeService = sn.CanBeService || sn.hasMod || (!sn.hasMod && (sn.hasGo)) //sn.hasMain || sn.hasMainGo || sn.hasMod || sn.hasGo
}

func (sn *ServiceNode) Print(prefix string) {
	stringToPrint := prefix + "|---" + sn.name
	if sn.CanBeService {
		stringToPrint += " $"
	}
	println(stringToPrint)
	for _, child := range sn.children {
		child.Print(prefix + "|   ")
	}
}

type ServiceFinder struct {
	sourcePath string
	nodes      []ServiceNode
	root       ServiceNode
	services   []string
}

func NewServiceFinder(sourcePath string) *ServiceFinder {
	return &ServiceFinder{
		sourcePath: sourcePath,
		nodes:      []ServiceNode{},
		root:       ServiceNode{},
		services:   []string{},
	}
}

func (finder *ServiceFinder) CreateNode(path string, parent *ServiceNode) *ServiceNode {
	node := ServiceNode{
		path:     path,
		parent:   parent,
		name:     filepath.Base(path),
		children: []ServiceNode{},
	}
	//if parent != nil {
	//	parent.children = append(parent.children, node)
	//}
	finder.nodes = append(finder.nodes, node)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.Name() == "go.mod" {
			node.hasMod = true
		} else if f.Name() == "main.go" {
			node.hasMainGo = true
			node.NGo++
		} else if strings.HasSuffix(f.Name(), ".go") {
			node.hasGo = true
			node.NGo++
			node.hasMain = HasMainFunction(filepath.Join(path, f.Name()))
		} else if f.IsDir() {
			// something else
			node.children = append(node.children, *finder.CreateNode(filepath.Join(path, f.Name()), &node))
		}
	}
	node.Update()
	return &node
}

func (finder *ServiceFinder) CreateRoot() {
	finder.root = *finder.CreateNode(finder.sourcePath, nil)
}

func (finder *ServiceFinder) findServicesRoot(currentNode *ServiceNode) (string, []string) {
	potentialServices := []ServiceNode{}
	for _, child := range currentNode.children {
		if child.CanBeService {
			potentialServices = append(potentialServices, child)
		}
	}
	if len(potentialServices) == 0 {
		if currentNode.CanBeService {
			return currentNode.path, []string{currentNode.path}
		} else {
			// raise error TODO
			log.Println("No services found in the directory", currentNode.path)
			return "", []string{} //, "No services found in the directory
		}
	} else if len(potentialServices) == 1 {
		return finder.findServicesRoot(&potentialServices[0])
	} else {
		if currentNode.NGo != 0 {
			log.Println("Multiple services found in the same directory", currentNode.path)
		}
		services := []string{}
		for _, service := range potentialServices {
			services = append(services, service.path)
		}
		return currentNode.path, services
	}
	//return "", nil
}

func (finder *ServiceFinder) GetServices() (string, []string) {
	return finder.findServicesRoot(&finder.root)
}

func HasMainFunction(path string) bool {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		if fn, isFn := decl.(*ast.FuncDecl); isFn {
			if fn.Name.Name == "main" {
				return true
			}
		}
	}

	return false
}
