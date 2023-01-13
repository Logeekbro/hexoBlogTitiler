package main

import "os"

/**
定义标题添加规则的抽象层
*/

type ITitler interface {
	// GetName 获取当前规则定义器名称
	GetName() string

	// IsTargetFile 判断文件是否为需要添加标题的文件
	IsTargetFile(filename string) bool

	// HasTitle 判断文件是否已有标题
	HasTitle(file *os.File) (bool, error)

	// GetTitle 获取标题
	GetTitle(filename string) string
}
