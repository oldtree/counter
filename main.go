// Counter project main.go
package main

import (
	"Counter/parser"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"os"

	"flag"

	"github.com/toolkits/file"
)

func main() {

	InputDir := flag.String("i", ".", "input dir")
	OutputDir := flag.String("o", ".", "output dir")
	flag.Parse()
	err := ParserDirFuncsList(*InputDir, *OutputDir)
	if err != nil {
		fmt.Println(err)
	}
}

func ParserDirFuncsList(dirPath string, outPutPath string) error {
	fmt.Println(dirPath, outPutPath)
	if !file.IsExist(dirPath) {
		return errors.New("input dir not exist")
	}
	if !file.IsExist(outPutPath) {
		return errors.New("output dir not exist")
	}

	walkfuncs := func(path string, info os.FileInfo, err error) error {
		if file.Ext(path) != ".go" {
			return nil
		}
		goFile, err := parser.ParseFile(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		parser.BuildFuncsTable(goFile)
		return nil
	}
	filepath.Walk(dirPath, filepath.WalkFunc(walkfuncs))
	var totallist = new(parser.FuncCounter)
	for _, value := range parser.FuncsTable {
		totallist.FuncsList = append(totallist.FuncsList, value.(*parser.FuncCounter).FuncsList...)
		totallist.FuncsNumber = totallist.FuncsNumber + value.(*parser.FuncCounter).FuncsNumber
	}
	totallist.PackagePath = dirPath
	relative := strings.Split(dirPath, "\\")
	var relativePath string
	if relative[len(relative)-1] == "" {
		relativePath = relative[len(relative)-2]
	} else {
		relativePath = relative[len(relative)-1]
	}
	data, _ := json.Marshal(totallist)
	fi, _ := os.Create(outPutPath + "/" + relativePath + ".json")
	fi.Write(data)
	fi.Sync()
	defer fi.Close()
	return nil
}
