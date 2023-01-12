package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const FileSep = "\\"
const ReadSize = 1024 * 100

func main() {
	var path string
	fmt.Println("Input your posts absolute path:")
	_, err := fmt.Scanln(&path)
	if err != nil {
		fmt.Println("Get path error:", err)
		return
	}
	readDir(path)
	fmt.Println("Add title complete!")
}

func readDir(path string) {
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Read Dir failed:%v\n", err)
		return
	}
	for _, info := range fileInfos {
		fullPath := path + FileSep + info.Name()
		if info.IsDir() {
			readDir(fullPath)
		} else if filepath.Ext(info.Name()) == ".md" {
			file, err := os.OpenFile(fullPath, os.O_RDWR, 777)
			if err != nil {
				fmt.Printf("Open file(path: %s) error: %s\n", fullPath, err)
				file.Close()
				continue
			}
			threeB := make([]byte, 3)
			_, err = file.Read(threeB)
			if err != nil {
				fmt.Printf("Read file(path: %s) error: %s\n", fullPath, err)
				file.Close()
				continue
			}
			// 以三个短横线开头则认为已经有标题
			if string(threeB) == "---" {
				file.Close()
				continue
			}
			// 将文件指针跳回文件开头位置
			_, err = file.Seek(0, 0)
			// 将原内容读取出来
			fileData := make([]byte, ReadSize)
			cnt, err := file.Read(fileData)
			if err != nil {
				fmt.Printf("Read fileData(path: %s) error: %s\n", fullPath, err)
				file.Close()
				continue
			} else if cnt > ReadSize {
				fmt.Printf("Can not read > 100MB file(path: %s)!\n", fullPath)
				file.Close()
				continue
			}
			// 将文件指针跳回文件开头位置
			_, err = file.Seek(0, 0)
			// 获取标题
			title := strings.TrimSuffix(info.Name(), ".md")
			formatData := "---\n" + "title: " + title + "\n---\n" + string(fileData[:cnt])
			// 重新将带开头的内容写入, 此处要将fileData截取为实际读到的长度
			_, err = file.WriteString(formatData)
			if err != nil {
				fmt.Printf("Write file(path: %s) error: %s\n", fullPath, err)
				file.Close()
				continue
			}
			file.Close()
			fmt.Printf("file %s add formatData success!\n", info.Name())
		}
	}
}
