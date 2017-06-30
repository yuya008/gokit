package command

import (
	"os"
	"fmt"
	buil "github.com/yuya008/gokit/builder"
	"os/user"
	"path"
	"gopkg.in/urfave/cli.v2"
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
	Authors = map[string]string{
		"ArthurYu": "yuya008@aliyun.com",
	}
	commands []*cli.Command
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
	gokitApp := cli.App{
		Name: ProjectName,
		Version: Version,
		Usage: "Go project build tools",
		Authors: authors(),
		Commands: commands,
	}
	return gokitApp.Run(os.Args)
}

func authors() []*cli.Author {
	var ret []*cli.Author
	for author, email := range Authors {
		ret = append(ret, &cli.Author{
			Name: author,
			Email: email,
		})
	}
	return ret
}

