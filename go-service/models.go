package main

type Object_ struct {
	IsInterface  bool   `json:"isInterface"`
	IsAnnotation bool   `json:"isAnnotation"`
	SimpleName   string `json:"simpleName"`
	FullName     string `json:"fullName"`
	FilePath     string `json:"filePath"`
	ServiceName  string `json:"serviceName"`
	Content      string `json:"content"`
}

type Executable_ struct {
	FullName    string `json:"fullName"`
	SimpleName  string `json:"simpleName"`
	ParentName  string `json:"parentName"`
	ServiceName string `json:"serviceName"`
	Content     string `json:"content"`
}
