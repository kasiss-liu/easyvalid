package pkg

import (
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

type parserEntry struct {
	filename   string
	pkg        *packages.Package
	syntax     *ast.File
	structures []StructInfo
}

//Parser 项目解析器
type Parser struct {
	pkg               *packages.Package
	entries           []*parserEntry
	entriesByFileName map[string]*parserEntry
	parserPackages    []*types.Package
	conf              packages.Config
}

//NewParser 构建并返回一个新的解析器
func NewParser() *Parser {
	var conf packages.Config
	conf.Mode = packages.LoadSyntax
	conf.Env = []string{"GO111MODULE=on"}

	return &Parser{
		parserPackages:    make([]*types.Package, 0),
		entriesByFileName: map[string]*parserEntry{},
		conf:              conf,
	}
}

//Parse 解析一个文件名或文件夹
func (p *Parser) Parse(path string) (err error) {

	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	var f os.FileInfo
	if f, err = os.Stat(path); err != nil {
		return nil
	}
	dir := filepath.Dir(path)
	files := []os.FileInfo{f}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".go" || strings.HasSuffix(file.Name(), "_test.go") {
			continue
		}
		fname := file.Name()
		fpath := filepath.Join(dir, fname)
		if _, ok := p.entriesByFileName[fpath]; ok {
			continue
		}

		//解析pkg
		var pkgs []*packages.Package
		pkgs, err = packages.Load(&p.conf, "file="+fpath)
		if err != nil {
			return
		}

		if len(pkgs) == 0 {
			continue
		}

		//一个项目内有多个package name 不合法。
		if len(pkgs) > 1 {
			names := make([]string, len(pkgs))
			for i, p := range pkgs {
				names[i] = p.Name
			}
			panic(fmt.Sprintf("file %s resolves to multiple packages: %s", fpath, strings.Join(names, ", ")))
		}

		pkg := pkgs[0]
		if len(pkg.Errors) > 0 {
			// err = pkg.Errors[0]
			// return
			return pkg.Errors[0]
		}

		if len(pkg.GoFiles) == 0 {
			continue
		}

		p.pkg = pkg

		for idx, f := range pkg.GoFiles {
			if _, ok := p.entriesByFileName[f]; ok {
				continue
			}
			entry := parserEntry{
				filename: f,
				pkg:      pkg,
				syntax:   pkg.Syntax[idx],
			}
			p.entries = append(p.entries, &entry)
			p.entriesByFileName[f] = &entry
		}

	}

	return
}

//Load 加载已经解析过的文件中的结构体
func (p *Parser) Load() error {
	for _, entry := range p.entries {
		nv := NewNodeVisitor()
		ast.Walk(nv, entry.syntax)
		entry.structures = nv.DeclaredStructures()
	}
	if p.GetPkg() == nil {
		return fmt.Errorf("no packages found")
	}
	if len(p.entries) == 0 {
		return fmt.Errorf("no structs need to be created")
	}
	return nil
}

//GetStructures 获取已经声明的结构体
func (p *Parser) GetStructures() []*Structure {
	structures := make(sortableStructures, 0)
	for _, entry := range p.entries {
		declaredStructures := entry.structures
		structures = p.packageStructures(entry.pkg.Types, entry.filename, declaredStructures, structures)
	}
	sort.Sort(structures)
	return structures
}

//GetStructure 根据需要的结构体名称 获取已经解析的结构体
func (p *Parser) GetStructure(stnames []string) []*Structure {
	structures, sts := make(sortableStructures, 0), make(sortableStructures, 0)
	for _, entry := range p.entries {
		declaredStructures := entry.structures
		sts = p.packageStructures(entry.pkg.Types, entry.filename, declaredStructures, structures)
	}
	for _, name := range stnames {
		for _, st := range sts {
			if st.Name == name {
				structures = append(structures, st)
				break
			}
		}
	}
	sort.Sort(structures)
	return structures
}

//GetImports 获取类中的package引用
func (p *Parser) GetImports() []*packages.Package {

	pkgs := make([]*packages.Package, 0)
	for _, entry := range p.entries {
		imps := entry.pkg.Imports
		for _, v := range imps {
			pkgs = append(pkgs, v)
		}
	}
	return pkgs
}

//GetPkg 获取被解析文件的包信息
func (p *Parser) GetPkg() *packages.Package {
	return p.pkg
}

//packageStructures 解析返回结构体信息
func (p *Parser) packageStructures(
	pkg *types.Package,
	fileName string,
	declaredStructures []StructInfo,
	structures []*Structure) []*Structure {
	scope := pkg.Scope()

	for _, info := range declaredStructures {
		obj := scope.Lookup(info.Name)
		if obj == nil {
			continue
		}

		typ, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}

		name := typ.Obj().Name()
		st, ok := typ.Underlying().(*types.Struct)
		if !ok {
			continue
		}

		if typ.Obj().Pkg() == nil {
			continue
		}

		elem := &Structure{
			Name:          name,
			Pkg:           pkg,
			QualifiedName: pkg.Path(),
			FileName:      fileName,
			Struct:        st,
			NamedType:     typ,
			Alias:         info.Alias,
		}

		structures = append(structures, elem)
	}

	return structures
}

//NodeVisitor 节点遍历器
type NodeVisitor struct {
	declaredStructures []StructInfo
}

//StructInfo 简单的结构体信息
type StructInfo struct {
	Name  string //结构体名称
	Alias string //注释中的结构体 方法别名 easyvalid:valid:(alias)
}

//NewNodeVisitor 构建并返回新的遍历器
func NewNodeVisitor() *NodeVisitor {
	return &NodeVisitor{
		declaredStructures: make([]StructInfo, 0),
	}
}

//DeclaredStructures 返回已经声明的结构体
func (nv *NodeVisitor) DeclaredStructures() []StructInfo {
	return nv.declaredStructures
}

//Visit 执行遍历
func (nv *NodeVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		if alias, ok := checkGenerate(n.Doc.Text(), generateFlag); ok {
			if _, ok := n.Type.(*ast.StructType); ok {
				nv.declaredStructures = append(nv.declaredStructures, StructInfo{n.Name.Name, alias})
			}
		}
	case *ast.GenDecl:
		for _, nc := range n.Specs {
			switch nct := nc.(type) {
			case *ast.TypeSpec:
				nct.Doc = n.Doc
			}
		}

	}

	return nv
}

//checkGenerate 校验是否可以生成内容 如果未声明结构体别名 默认用this
//func (this estruct) Valid() {}
func checkGenerate(comment, flag string) (string, bool) {
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, flag) {
			alias := "this"
			if ws := strings.Split(line, ":"); len(ws) > 2 {
				alias = ws[2]
			}
			return alias, true
		}
	}
	return "", false
}

//Structure 结构体
type Structure struct {
	Name          string
	Alias         string
	QualifiedName string
	FileName      string
	File          *ast.File
	Pkg           *types.Package
	Struct        *types.Struct
	NamedType     *types.Named
}

//GetTag 查找tag中的指定key
//方法同 reflect.Elem().Field(i).Tag(i).Get(key)
func (st *Structure) GetTag(index int, key string) (value string, ok bool) {

	if st.Struct.NumFields() <= index {
		ok = false
		return
	}
	//获取tag内容
	tag := st.Struct.Tag(index)

	for tag != "" {
		// 跳过前驱空格符.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// 扫描引号字符 匹配内容
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}

//sortableStructures 对获取到的结构体进行排序
type sortableStructures []*Structure

func (s sortableStructures) Len() int {
	return len(s)
}

func (s sortableStructures) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortableStructures) Less(i, j int) bool {
	return strings.Compare(s[i].Name, s[j].Name) == -1
}
