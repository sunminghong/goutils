package logger

import (
  "fmt"

	)

func Debug(msg ...interface{}){
  fmt.Print("DEBUG: ")
  fmt.Println(msg...)
}


func Info(msg ...interface{}){
  fmt.Print("INFO: ")
  fmt.Println(msg...)
}


func Error(msg ...interface{}){
  fmt.Print("ERROR: ")
  fmt.Println(msg...)
}

func Warn(msg ...interface{}){
  fmt.Print("WARN: ")
  fmt.Println(msg...)
}

func Fatal(msg ...interface{}){
  fmt.Print("FATAL: ")
  fmt.Println(msg...)
  panic(msg)
}
