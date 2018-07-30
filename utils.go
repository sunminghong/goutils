package utils

import (
	"path/filepath"
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


