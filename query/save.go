package query

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
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

func Write(to io.Writer, table *Table) error {
	save := gob.NewEncoder(to)
	return save.Encode(table)
}

func Read(from io.Reader, table *Table) error {
	load := gob.NewDecoder(from)
	return load.Decode(&table)
}

func Save(path string, table *Table) (wtf error) {

	defer func() {
		wtf = catch(recover(), "Save")
	}()

	file, err := os.Create(path)
	assert(err)
	defer file.Close()

	zip := gzip.NewWriter(file)
	defer zip.Close()

	save := gob.NewEncoder(zip)
	assert(save.Encode(table))

	return
}

func Load(path string) (table *Table, wtf error) {

	defer func() {
		wtf = catch(recover(), "Load")
	}()

	file, err := os.Open(path)
	assert(err)
	defer file.Close()

	zip, zerr := gzip.NewReader(file)
	assert(zerr)
	defer zip.Close()

	load := gob.NewDecoder(zip)
	assert(load.Decode(&table))

	return
}
