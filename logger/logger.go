package logger

import (
  "fmt"
  "time"

	)

func getTime() string {
	f := "2006-01-02 15:04:05.000@"
  return time.Now().Local().Format(f)
}

func Debug(msg ...interface{}){
  return
  fmt.Print(getTime())
  fmt.Print("DEBUG: ")
  fmt.Println(msg...)
}


func Dot(f string){
  fmt.Print(f)
}

func Info(msg ...interface{}){
  fmt.Print(getTime())
  fmt.Print("INFO: ")
  fmt.Println(msg...)
}


func Error(msg ...interface{}){

  /*
	pc, file, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fmt.Println("【pkgame】", fmt.Sprintf("func = %s,file = %s,line = %d,ok = %v ,val = %v", f.Name(), file, line, ok, txt))

  */

  fmt.Print(getTime())
  fmt.Print("ERROR: ")
  fmt.Println(msg...)
}

func Warn(msg ...interface{}){
  fmt.Print(getTime())
  fmt.Print("WARN: ")
  fmt.Println(msg...)
}

func Fatal(msg ...interface{}){
  fmt.Print(getTime())
  fmt.Print("FATAL: ")
  fmt.Println(msg...)
  panic(fmt.Sprint(msg...))
}
