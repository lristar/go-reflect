package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name        string
	Age         int
	Sex         uint
	Description string
}

func main() {
	var a = Student{Name:"张三",Age:20,Sex:1,Description:"hahaha"}
	// 获取类型
	//elems := reflect.TypeOf(&a).Elem()
	//for i := 0; i < elems.NumField(); i++ {
	//	field := elems.Field(i)
	//}
	// 获取数值
	elems := reflect.ValueOf(&a).Elem()
	for i :=0;i<elems.NumField();i++{
		fmt.Println(elems.Field(i))
	}


}
