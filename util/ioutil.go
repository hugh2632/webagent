package util

import (
	"errors"
	"os"
	"reflect"
	"regexp"
)

//利用正则表达式压缩字符串，去除空格或制表符
func Trim(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, " ")
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contains(target interface{}, obj interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}

//IsExist  判断文件夹/文件是否存在  存在返回 true
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func EnsurePath(path string) error {
	if !IsExist(path) {
		return CreateDir(path)
	}
	return nil
}

func EnsureFile(file string) (*os.File, error) {
	if !IsExist(file) {
		return os.Create(file)
	} else {
		return os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	}
}

//CreateDir  文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Chmod(path, os.ModePerm)
	return err
}
