package handler

import (
	"fmt"
	"github.com/d5/tengo/script"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
	"io/ioutil"
	"log"
	"rsps/entity"
	"rsps/model"
)

type ObjectObserver struct {
	Function       interface{}
	CompiledScript *script.Compiled
}


var ObjectObservers = make(map[int]interface{})

var scriptsDir = "./scripts2"

//func LoadScripts() error {
//	files, err := ioutil.ReadDir(scriptsDir)
//	if err != nil {
//		panic(err)
//	}
//
//	importModules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
//	for _, file := range files {
//		var c *script.Compiled
//		code, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", scriptsDir, file.Name()))
//		if err != nil {
//			panic(err)
//		}
//		s := script.New(code)
//		s.SetImports(importModules)
//
//		bindObjectFirstClick := &objects.UserFunction{
//			Value: func(args ...objects.Object) (ret objects.Object, err error) {
//				objectId := args[0]
//				f := args[1]
//				ob := &ObjectObserver{
//					Function:       f,
//					CompiledScript: c,
//				}
//
//				oid, _ := strconv.Atoi(objectId.String())
//				ObjectObservers[oid] = ob
//				return nil, nil
//			},
//		}
//
//		execute := &objects.UserFunction{
//			Value: func(args ...objects.Object) (ret objects.Object, err error) {
//				log.Printf("called empty execute")
//				return nil, nil
//			},
//		}
//
//		s.Add("state", "bind")
//		s.Add("execute", execute)
//		s.Add("bindObjectFirstClick", bindObjectFirstClick)
//		s.Add("object", nil)
//		s.Add("player", nil)
//		s.Add("setObject", SetObject())
//
//		c, err = s.Compile()
//		if err != nil {
//			panic(err)
//		}
//
//		if err := c.Run(); err != nil {
//			panic(err)
//		}
//
//		c.Set("state", "execute")
//	}
//
//	return nil
//}
//
func SetObject(id, x, y int) {
	entity.WorldProvider().AddWorldObject(id, &model.Position{
		X: uint16(x),
		Y: uint16(y),
		Z: 0,
	})
}


func LoadScripts() {
	files, err := ioutil.ReadDir(scriptsDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", scriptsDir, file.Name()))
		if err != nil {
			log.Println("Error reading script file for object action:", err)
		}

		e := WorldModule()
		e.Define("printf", fmt.Printf)
		_, err = vm.Execute(e, nil, string(data))
		if err != nil {
			log.Println("error binding:", err)
		}
	}
}

func Run() {
	//w := WorldModule()
}

func WorldModule() *env.Env {
	e := env.NewEnv()
	_ = e.Define("bind", map[string]interface{}{
		"object": func(id int, f interface{}) {
			ObjectObservers[id] = f
		},
	})

	_ = e.Define("world", map[string]interface{}{
		"setObject": SetObject,
	})

	return e
}
