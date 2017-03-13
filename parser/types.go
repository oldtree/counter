package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type GoFile struct {
	Package    string
	Path       string
	Structs    []*GoStruct
	Interfaces []*GoInterface
	Imports    []*GoImport
	FuncsList  []*GoFuncs
}

type FuncCounter struct {
	PackagePath string     `json:"package"`
	FuncsList   []*GoFuncs `json:"funclist"`
	FuncsNumber int        `json:"number"`
}

var FuncsTable = make(map[string]interface{})

func BuildFuncsTable(g *GoFile) {
	if g == nil {
		return
	}
	if _, ok := FuncsTable[g.Path]; ok {
		return
	}
	var cos = new(FuncCounter)
	cos.FuncsList = append(cos.FuncsList, g.FuncsList...)
	cos.FuncsNumber = len(cos.FuncsList)
	FuncsTable[g.Path] = cos
	return
}

func (g *GoFile) ImportPath() (string, error) {
	importPath, err := filepath.Abs(g.Path)
	if err != nil {
		return "", err
	}

	importPath = strings.Replace(importPath, "\\", "/", -1)

	goPath := strings.Replace(os.Getenv("GOPATH"), "\\", "/", -1)
	importPath = strings.TrimPrefix(importPath, goPath)
	importPath = strings.TrimPrefix(importPath, "/src/")

	importPath = strings.TrimSuffix(importPath, filepath.Base(importPath))
	importPath = strings.TrimSuffix(importPath, "/")

	return importPath, nil
}

type GoImport struct {
	File *GoFile
	Name string
	Path string
}

type GoInterface struct {
	File    *GoFile
	Name    string
	Methods []*GoMethod
}

type GoMethod struct {
	Name    string
	Params  []*GoType
	Results []*GoType
}

type GoFuncs struct {
	Name    string `json:"funcname"`
	Reciver string `json:"reciver"`
}

func (g *GoFuncs) FuncsSpesKey() string {
	if g.Reciver == "" {
		return g.Name
	} else {
		return g.Reciver + "." + g.Name
	}
}

type GoType struct {
	Name  string
	Type  string
	Inner []*GoType
}

type GoStruct struct {
	File   *GoFile
	Name   string
	Fields []*GoField
}

type GoField struct {
	Struct *GoStruct
	Name   string
	Type   string
	Tag    *GoTag
}

type GoTag struct {
	Field *GoField
	Value string
}

func (g *GoTag) Get(key string) string {
	tag := strings.Replace(g.Value, "`", "", -1)
	return reflect.StructTag(tag).Get(key)
}
func (g *GoImport) Prefix() string {
	if g.Name != "" {
		return g.Name
	}

	path := strings.Trim(g.Path, "\"")
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return path
	}

	return path[lastSlash+1:]
}
