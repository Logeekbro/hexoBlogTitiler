package main

import (
	"os"
	"path/filepath"
	"strings"
)

/**
定义md文件标题添加规则
*/

type MdTitler struct {
	Name string
}

func NewMdTitler() ITitler {
	return &MdTitler{
		Name: "MdTitler",
	}
}

func (t *MdTitler) GetName() string {
	return t.Name
}

func (t *MdTitler) IsTargetFile(filename string) bool {
	return filepath.Ext(filename) == ".md"
}

func (t *MdTitler) HasTitle(file *os.File) (bool, error) {
	// 读取文件开头的3个字节
	threeB := make([]byte, 3)
	_, err := file.Read(threeB)
	if err != nil {
		return true, err
	}
	// 以三个短横线开头则认为已经有标题
	if string(threeB) == "---" {
		return true, nil
	}
	return false, nil
}

func (t *MdTitler) GetTitle(filename string) string {
	return "---\n" + "title: " + strings.TrimSuffix(filename, ".md") + "\n---\n"
}
