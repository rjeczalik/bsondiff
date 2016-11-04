package bsondiff

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func Diff(old, new map[string]interface{}, out *map[string]interface{}) error {
	type elem struct {
		keys []string
		old  map[string]interface{}
		new  map[string]interface{}
	}

	var (
		cur    elem
		stack  = []elem{{old: old, new: new}}
		set    = make(map[string]interface{})
		unset  = make(map[string]interface{})
		update = make(map[string]interface{})
	)

	for len(stack) != 0 {
		cur, stack = stack[0], stack[1:]

		oldKeys := keys(cur.old)

		for _, oldKey := range oldKeys {
			v, ok := cur.new[oldKey]
			if !ok {
				set(unset, append(cur.keys, oldKey), oldKeys[oldKey])
				continue
			}
		}
	}

	return nil
}

func keys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func set(m map[string]interface{}, k []string, v interface{}) {
	for _, k := range k[:len(k)-1] {
		mm, ok = m[k].(map[string]interface{})
		if !ok {
			mm = make(map[string]interface{})
			m[k] = mm
		}

		m = mm
	}

	m[len(k)-1] = v
}

type Program struct {
	JSON   bool
	Stdout io.Writer
}

func (p *Program) Run(f *flag.FlagSet, args []string) error {
	f.BoolVar(&p.JSON, "json", false, "")

	if err := f.Parse(args); err != nil {
		return err
	}

	if p.JSON {
		cols, err := p.readCollections(f.Args()...)
		if err != nil {
			return err
		}

		enc := json.NewEncoder(p.Stdout)
		enc.SetIndent("", "\t")

		return enc.Encode(cols)
	}

	if f.NArg() != 2 {
		return errors.New("usage: bsondiff <old dump dir> <new dump dir>")
	}

	oldCols, err := p.readCollectionDir(f.Arg(0))
	if err != nil {
		return fmt.Errorf("reading %q: %s", f.Arg(0), err)
	}

	newCols, err := p.readCollectionDir(f.Arg(1))
	if err != nil {
		return fmt.Errorf("reading %q: %s", f.Arg(1), err)
	}

	_, _ = oldCols, newCols

	return nil
}

func (p *Program) readFiles(path string) (files []string, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		fis, err := ioutil.ReadDir(fi.Name())
		if err != nil {
			return nil, err
		}

		for _, fi := range fis {
			files = append(files, fi.Name())
		}
	} else {
		files = append(files, path)
	}

	return files, nil
}

func (p *Program) readCollections(files ...string) (map[string]interface{}, error) {
	cols := make(map[string]interface{})

	for _, file := range files {
		name := filepath.Base(file)

		i := strings.Index(name, ".bson")
		if i == -1 {
			continue
		}

		name = name[:i]

		p, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		var v interface{}

		if err := bson.Unmarshal(p, &v); err != nil {
			return nil, fmt.Errorf("reading %q: %s", file, err)
		}

		cols[name] = v
	}

	return cols, nil
}

func (p *Program) readCollectionDir(dir string) (map[string]interface{}, error) {
	files, err := p.readFiles(dir)
	if err != nil {
		return nil, err
	}

	return p.readCollections(files...)
}
