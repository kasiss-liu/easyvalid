package pkg

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

type A struct{}

var a A

func TestAAA(t *testing.T) {
	t.Log(&a == &A{})
}

func TestParse(t *testing.T) {
	parser := NewParser()
	err := parser.Parse("../_examples/Animal/Dog.go")
	if err != nil {
		t.Error(err.Error())
	}

	parser.Load()

	pkgs := parser.GetImports()
	t.Log(pkgs)
	for k, v := range pkgs {
		fmt.Println(k, v)
	}

	for _, st := range parser.GetStructures() {
		fmt.Println("filename", st.FileName)
		fmt.Println("name", st.Name)
		fmt.Println("tag", st.Struct.Tag(0), st.Struct.Tag(1))
	}
}

func TestStructures(t *testing.T) {
	st := &Structure{}
	tt := reflect.TypeOf(st)
	valid := tt.Elem().Field(0).Tag.Get("valid")
	fmt.Println("validtag", valid)
}

func TestParse2(t *testing.T) {
	dir := "../../downloader/"
	p := NewParser()
	walk(p, dir)
	p.Load()
	for _, st := range p.GetStructures() {
		fmt.Println("filename", st.FileName)
		fmt.Println("name", st.Name)
		fmt.Println("tag", st.Struct.Tag(0), st.Struct.Tag(1))
	}
}

func walk(p *Parser, dir string) (err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") || strings.HasPrefix(file.Name(), "_") {
			continue
		}

		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			walk(p, path)
			continue
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			continue
		}

		err = p.Parse(path)
		if err != nil {
			continue
		}
	}
	return
}

func TestExample(t *testing.T) {

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedImports}
	pkgs, err := packages.Load(cfg, "../_examples/Animal/Animal.go")
	if err != nil {
		t.Errorf("load: %v\n", err)
		return
	}
	if packages.PrintErrors(pkgs) > 0 {
		t.Errorf("pkg errors: %v\n", err)
		return
	}
	for _, pkg := range pkgs {
		fmt.Println(1, pkg.ID, pkg.GoFiles, pkg.Imports)
	}
}
