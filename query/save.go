package query

import (
	"encoding/gob"
	"fmt"
	"os"
)

func catch(r interface{}, decorate string) error {
	if r != nil {
		if _, ok := r.(error); ok {
			return fmt.Errorf("%s: %s", decorate, r.(error).Error())
		}
		return fmt.Errorf("%s: %v", decorate, r)
	}
	return nil
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func Save(path string, tables ...*Table) (wtf error) {

	defer func() {
		wtf = catch(recover(), "Save")
	}()

	file, err := os.Create(path)
	assert(err)
	defer file.Close()

	save := gob.NewEncoder(file)
	assert(save.Encode(tables))

	return
}

func Load(path string) (tables []*Table, wtf error) {

	defer func() {
		wtf = catch(recover(), "Load")
	}()

	file, err := os.Open(path)
	assert(err)
	defer file.Close()

	load := gob.NewDecoder(file)
	assert(load.Decode(&tables))

	return
}
