// #############################################################################
// # File: fileutil.go                                                         #
// # Project: fileutil                                                         #
// # Created Date: 2024/10/08 17:54:07                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2024/11/11 12:43:29                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// #                                                                           #
// #############################################################################
package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 获取程序运行路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func Abs(dir string) string {
	dir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	return dir
}

// 判断文件路径是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 判断是否是文件夹
func IsDir(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}
	return false
}

// 如文件不存在则创建
func CreateFileIfNecessary(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		if file, err := os.Create(path); err == nil {
			file.Close()
		}
	}
	exist := IsExist(path)
	return exist
}

// 如目录不存在则创建
func MkdirIfNecessary(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			// os.Chmod(path, 0777)
			return err
		}
	}
	return nil
}
