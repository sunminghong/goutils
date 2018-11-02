package utils

import (
	"path/filepath"
  "io/ioutil"
	"os"
	"fmt"
	"strings"
  "strconv"
	"crypto/sha1"

  "github.com/satori/go.uuid"

)

func GetSha(source string) string {
	h:= sha1.New()
	h.Write([]byte(source))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}


 func Toint(a string) int {
   //if i, err := strconv.Atoi(a); err != nil {
   //return i
   //}
   //return 0
   i, _ := strconv.Atoi(a)
   return i
 }

func Uuid() string {
  ui,_ :=uuid.NewV4()
  return ui.String()
}

/*
获取程序运行路径
*/
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("This is an error")
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	//PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, fi.Name())
		}
	}

	return files, nil
}


