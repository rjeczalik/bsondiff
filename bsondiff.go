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
	"strings"

	"github.com/kylelemons/godebug/pretty"
	"gopkg.in/mgo.v2/bson"
)

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

	fmt.Fprintln(p.Stdout, pretty.Compare(oldCols, newCols))

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
