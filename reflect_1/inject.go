package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const TAGINJECT = "inject"

type Service struct {
	Name string
}

func (s *Service) Hello() {
	fmt.Println("hello")
}

type Controller struct {
	CService *Service `inject:"Service"`
}

type Factory struct {
	CController *Controller `inject:"Controller"`
}

type Container struct {
	sync.Mutex
	singleton map[string]interface{}
}

var GFactory = &Container{singleton: map[string]interface{}{}}

// 写入方法

func (c *Container) SetSingleton(v interface{}) {
	c.Lock()
	c.singleton[GetStructName(v)] = v
	c.Unlock()
}

func (c *Container) GetSingleton(key string) (interface{}, error) {
	if v, ok := c.singleton[key]; ok {
		return v, nil
	}
	return nil, errors.New("没有对应的对象")
}

// 注入实例

func (c *Container) Entry(v interface{}) error {
	if err := c.EntryValue(reflect.ValueOf(v)); err != nil {
		return err
	}
	return nil

}

func (c *Container) EntryValue(v reflect.Value) error {
	// 先判断这个是不是指针，不是就退出
	if v.Kind() != reflect.Ptr {
		return errors.New("不是指针退出")
	}
	types, values := v.Type().Elem(), v.Elem()
	length := types.NumField()
	for i := 0; i < length; i++ {
		if !values.Field(i).CanSet() {
			continue
		}
		if types.Field(i).Anonymous {
			continue
		} else {
			// 首先要确认这个参数是为空的，没有被赋值的
			if values.Field(i).IsZero() {
				// 获取tag
				tag := types.Field(i).Tag.Get(TAGINJECT)
				fun, err := c.GetSingleton(tag)
				if err != nil {
					return err
				}
				if err := c.EntryValue(reflect.ValueOf(fun)); err != nil {
					return err
				}
				values.Field(i).Set(reflect.ValueOf(fun))
				fmt.Println("注入了", types.Field(i), "-", reflect.ValueOf(fun).Type())
			} else {
				fmt.Println("非空", values.Field(i))
			}
		}
	}
	return nil
}

func Override(root interface{}, v ...interface{}) error {
	GFactory.SetSingleton(root)
	for _, values := range v {
		GFactory.SetSingleton(values)
	}
	err := GFactory.Entry(root)
	if err != nil {
		return err
	}
	return nil
}
func GetStructName(v interface{}) string {
	names := reflect.TypeOf(v).String()
	name := strings.Split(names, ".")
	return name[len(name)-1]
}

func Init() {
	ctlFactory := &Factory{}
	err := Override(ctlFactory, &Service{Name: "张三"}, &Controller{})
	if err != nil {
		return
	}
}

func main() {
	Init()
	fac, err := GFactory.GetSingleton("Factory")
	if err != nil {
		fmt.Println("失败")
		return
	}
	fu := fac.(*Factory)
	fu.CController.CService.Hello()
}
