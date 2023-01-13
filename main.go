package main

import (
	"fmt"
	"io/fs"
	"os"
)

const (
	FileSep  = string(os.PathSeparator) // 文件分隔符
	ReadSize = 1024 * 100               // 读取文件的大小限制，单位：byte
)

// ErrorFiles 添加标题过程中出错的文件
var ErrorFiles = make([]string, 0)
var titler ITitler

func main() {
	var path string
	fmt.Println("Input your posts absolute path:")
	_, err := fmt.Scanln(&path)
	if err != nil {
		fmt.Println("Get path error:", err)
		return
	}
	// 设置命名规则
	titler = NewMdTitler()
	fmt.Printf("using %s...\n", titler.GetName())
	// 读取并处理文件
	handleDir(path)
	// 处理完成
	fmt.Println("Add title complete!")
	// 打印出错的文件
	if len(ErrorFiles) > 0 {
		fmt.Printf("%d files add title failed:\n", len(ErrorFiles))
		for index := range ErrorFiles {
			fmt.Println(ErrorFiles[index])
		}
	}
}

func handleDir(path string) {
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Read Dir failed:%v\n", err)
		return
	}
	for _, fileInfo := range fileInfos {
		fullPath := path + FileSep + fileInfo.Name()
		if fileInfo.IsDir() {
			handleDir(fullPath)
		} else if titler.IsTargetFile(fileInfo.Name()) {
			addTitle(fullPath, fileInfo)
		}
	}
}

func addTitle(fullPath string, fileInfo fs.DirEntry) {
	file, err := os.OpenFile(fullPath, os.O_RDWR, 777)
	defer file.Close()
	if err != nil {
		fmt.Printf("Open file(path: %s) error: %s\n", fullPath, err)
		ErrorFiles = append(ErrorFiles, fullPath)
		return
	}
	// 检测文件是否存在标题
	if hasTitle, err := titler.HasTitle(file); err != nil || hasTitle {
		if err != nil {
			fmt.Printf("Check has title error(path: %s): %s\n", fullPath, err)
			ErrorFiles = append(ErrorFiles, fullPath)
		}
		return
	}
	// 将文件指针跳回文件开头位置
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		ErrorFiles = append(ErrorFiles, fullPath)
		return
	}
	// 将原内容读取出来
	fileData := make([]byte, ReadSize)
	cnt, err := file.Read(fileData)
	if err != nil {
		fmt.Printf("Read fileData(path: %s) error: %s\n", fullPath, err)
		ErrorFiles = append(ErrorFiles, fullPath)
		return
	} else if cnt > ReadSize {
		fmt.Printf("Can not read > 100MB file!(path: %s)\n", fullPath)
		ErrorFiles = append(ErrorFiles, fullPath)
		return
	}
	// 将文件指针跳回文件开头位置
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		ErrorFiles = append(ErrorFiles, fullPath)
		return
	}
	// 重新将带标题的内容写入, 此处要将fileData截取为实际读到的长度
	formatData := titler.GetTitle(fileInfo.Name()) + string(fileData[:cnt])
	_, err = file.WriteString(formatData)
	if err != nil {
		fmt.Printf("Write file(path: %s) error: %s\n", fullPath, err)
		ErrorFiles = append(ErrorFiles, fullPath)
		return
	}
	fmt.Printf("file %s add title success!\n", fileInfo.Name())
}
