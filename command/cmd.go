package command

import (
	"os"
	"fmt"
	buil "github.com/yuya008/gokit/builder"
	"os/user"
	"path"
)

const (
	ProjectName = "gokit"
	Version = "0.0.1"
	GoKitWorksSapce = ".gokit"
	dirMode = 0755
)

var (
	builder *buil.Builder
	pwd string
)

type Command interface {
	ParseArgs([]string) (bool, error)
	Run() error
	Usage() string
}

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gokitPath := path.Join(u.HomeDir, GoKitWorksSapce)
	os.MkdirAll(gokitPath, dirMode)
	if builder, err = buil.NewBuilder(gokitPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if pwd, err = os.Getwd(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Execute() error {
	if len(os.Args) <= 1 {
		usage()
	}
	args := os.Args[1:]
	var command Command
	switch args[0] {
	case "version":
		printVersion()
	case "build":
		command = &CommandBuild{}
	case "run":
	
	case "test":
	
	case "import":
		command = &CommandImport{}
	default:
		usage()
	}
	if ok, err := command.ParseArgs(os.Args[2:]); !ok {
		fmt.Print(command.Usage())
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	if err := command.Run(); err != nil {
		fmt.Println(err)
	}
	return nil
}

func usage() {
	fmt.Println(`-----
gokit Go 工程构建工具链
-----
Usgae:
     gokit command [args]
Command:
     build       编译构建包
     run         编译构建包，并运行
     test        测试包
     import      导入包
     version     工具版本信息`)
	os.Exit(1)
}

func printVersion() {
	fmt.Printf("%s v%s\n", ProjectName, Version)
	os.Exit(0)
}