package cmd

import (
	"os"
	"fmt"
	"path"
	"flag"
	"log"
)

var (
	programName string
	vm VM
)

type VM struct {
	classPath string
}

func init() {
	programName = path.Base(os.Args[0])
}

func Run() error {
	if err := parseArgs(); err != nil {
		return err
	}
	log.Println(vm)
	return nil
}

func parseArgs() error {
	flagSet := flag.NewFlagSet(programName, flag.ExitOnError)
	flagSet.Usage = usage
	flagSet.StringVar(&vm.classPath, "-classpath", "", "")
	if err := flagSet.Parse(os.Args); err != nil {
		return err
	}
	return nil
}

func usage() {
	fmt.Printf(`用法: %s [-options] class [args...] (执行类)
或  %s [-options] -jar jarfile [args...] (执行 jar 文件)
其中选项包括:
	-cp <目录和 zip/jar 文件的类搜索路径>
	-classpath <目录和 zip/jar 文件的类搜索路径>
		用 : 分隔的目录, JAR 档案
		和 ZIP 档案列表, 用于搜索类文件。
	-version     输出产品版本并退出
	-? -help     输出此帮助消息
	-D<名称>=<值> 设置系统属性
`, programName, programName)
	os.Exit(1)
}