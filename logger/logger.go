package logger

import (
  "fmt"

	)

func Debug(msg ...interface{}){
  fmt.Println(msg...)
}


func Info(msg ...interface{}){
  fmt.Println(msg...)
}


func Error(msg ...interface{}){
  fmt.Println(msg...)
}

func Warn(msg ...interface{}){
  fmt.Println(msg...)
}
