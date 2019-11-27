package handler

import (
	"fmt"
	"github.com/d5/tengo/script"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
	"io/ioutil"
	"log"
	"os"
	"rsps/entity"
	"rsps/model"
)

type ObjectObserver struct {
	Function       interface{}
	CompiledScript *script.Compiled
}

var ObjectObservers = make(map[int]interface{})
var CommandObservers = make(map[string]interface{})

var scriptsDir = "./scripts"

func SetObject(id, x, y, face, typ int) {
	entity.WorldProvider().AddWorldObject(id,
		&model.Position{
			X: uint16(x),
			Y: uint16(y),
			Z: 0,
		},
		face,
		typ)
}

func GetObject(x, y int) model.WorldObjectInterface {
	return entity.WorldProvider().GetWorldObject(&model.Position{X: uint16(x), Y: uint16(y)})
}

func LoadScripts() {
	files, err := ioutil.ReadDir(scriptsDir)
	if err != nil {
		panic(err)
	}

	parseScripts(scriptsDir, files)
}

func parseScripts(directory string, files []os.FileInfo) {
	for _, file := range files {
		if file.IsDir() {
			dir, err := ioutil.ReadDir(scriptsDir + "/" + file.Name())
			if err != nil {
				continue
			}
			parseScripts(directory + "/" + file.Name(), dir)
			continue
		}
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", directory, file.Name()))
		if err != nil {
			log.Println("Error reading script file for object action:", err)
		}

		e := WorldModule()
		e.Define("printf", fmt.Printf)
		_, err = vm.Execute(e, &vm.Options{Debug: true}, string(data))
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
			log.Printf("bound obj %d", id)
			ObjectObservers[id] = f
		},
		"command": func(command string, f interface{}) {
			log.Printf("bound command %s", command)
			CommandObservers[command] = f
		},
	})

	_ = e.Define("world", map[string]interface{}{
		"setObject": SetObject,
		"getObject": GetObject,
	})

	_ = e.Define("NewPosition", model.NewPosition)

	return e
}
