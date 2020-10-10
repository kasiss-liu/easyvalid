package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kasiss-liu/easyvalid/pkg"
)

var (
	ppath       = flag.String("p", "", "package_path: a filepath or dir which package set")
	outname     = flag.String("o", "", "output filename")
	structnames = flag.String("structs", "", "struct valid need be generate. example: struct1,struct2")
	note        = flag.String("note", "", "note info to attach")
	buildTags   = flag.String("build_tags", "", "build flags to attach")
	help        = flag.Bool("h", false, "show help info")
	exported    = flag.Bool("export", false, "attr valid funcs if exported")
	force       = flag.Bool("f", false, "force rm file suffix with _easyvalid.go")
	verson      = flag.Bool("v", false, "show version")
)

func osExit(err error, code int) {
	fmt.Println(err)
	os.Exit(code)
}

func osExitWithHelp(err error, code int) {
	fmt.Println(err)
	flag.Usage()
	os.Exit(code)
}

func getFilename(dir string) (string, error) {
	//判断是否是文件夹
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return dir, nil
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", nil
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".go") {
			return dir + string(filepath.Separator) + file.Name(), nil
		}
	}
	return "", errors.New("no go files in path " + dir)
}

//生成输出文件名
func getOutputFilename(basedir, pkgname, outputname string) string {
	if outputname == "" {
		outputname = pkgname + "_easyvalid.go"
	}
	if !strings.HasSuffix(outputname, ".go") {
		outputname += ".go"
	}
	outputname = basedir + string(filepath.Separator) + outputname
	return outputname

}

func cleanOldEasyvalidFiles(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_easyvalid.go") {
			err := os.Remove(dir + string(filepath.Separator) + file.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *verson {
		osExit(errors.New(pkg.Version), 0)
	}

	//验证参数
	//验证输入文件
	if *ppath == "" {
		osExitWithHelp(errors.New("path is empty"), 1)
	}
	fdir, err := filepath.Abs(filepath.Dir(*ppath))
	if err != nil {
		osExit(err, 1)
	}
	//获取文件名
	fname, err := getFilename(*ppath)
	if err != nil {
		osExitWithHelp(err, 1)
	}
	//解析structs
	var stnames []string
	if strings.TrimSpace(*structnames) != "" {
		stnames = strings.Split(strings.TrimSpace(*structnames), ",")
	}
	*note = strings.TrimSpace(*note)
	*buildTags = strings.TrimSpace(*buildTags)
	*outname = strings.TrimSpace(*outname)

	//删除文件
	if *force {
		err = cleanOldEasyvalidFiles(fdir)
		if err != nil {
			osExit(err, 1)
		}
	}

	//解析文件
	parser := pkg.NewParser()
	err = parser.Parse(fname)

	if err != nil {
		osExit(err, 1)
	}
	err = parser.Load()
	if err != nil {
		osExit(err, 1)
	}

	//如果传入了指定的structure名
	var structs []*pkg.Structure
	if len(stnames) > 0 {
		structs = parser.GetStructure(stnames)
		//验证是否匹配
		stmap := make(map[string]bool)
		for _, st := range structs {
			stmap[st.Name] = true
		}
		for _, name := range stnames {
			if _, ok := stmap[name]; !ok {
				osExit(errors.New("struct "+name+" not found"), 1)
			}
		}
	} else {
		structs = parser.GetStructures()
	}
	//是否需要生成
	if len(structs) == 0 {
		osExit(errors.New("no easyvalid struct found"), 1)
	}

	gen := pkg.NewGenerator(parser.GetPkg(), structs, parser.GetImports())
	gen.SetNote(*note)
	gen.SetBuildTags(*buildTags)
	gen.SetValidFuncExport(*exported)

	err = gen.Run()
	if err != nil {
		osExit(err, 1)
	}
	buf := gen.GetPrintBuf()

	//生成输出文件名
	*outname = getOutputFilename(fdir, parser.GetPkg().Name, *outname)

	file, err := os.Create(*outname)
	if err != nil {
		osExit(err, 1)
	}

	_, err = fmt.Fprintln(file, buf.String())
	if err != nil {
		osExit(err, 1)
	}
	fmt.Println("github.com/kasiss-liu/easyvalid@" + pkg.Version)
	fmt.Println("file:", *outname)
	fmt.Println("generated success.")
}
