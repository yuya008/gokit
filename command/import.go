package command

import (
	"fmt"
	"os"
)

type CommandImport struct {
	packageName string
	source string
	insecure bool
}

func (ci *CommandImport) ParseArgs(args []string) (ok bool, error error) {
	var option string
	if len(args) < 1 {
		return false, nil
	}

	defer func() {
		if err := recover(); err != nil {
			ok = false
			error = fmt.Errorf("'%s' option value invalid", option)
		}
	}()

	ci.packageName = args[len(args) - 1]
	args = args[:len(args) - 1]
	for i := 0; i < len(args); i++ {
		option = args[i]
		switch option {
		case "--source":
			i++
			ci.source = args[i]
		case "--insecure":
			ci.insecure = true
		case "--help", "-h":
			return false, nil
		default:
			return false, fmt.Errorf("%s invalid option", args[i])
		}
	}
	return true, nil
}

func (ci *CommandImport) Run() error {
	if ci.source == "" {
		builder.Packager.Pull(ci.packageName, ci.insecure)
	} else {
		if f, err := os.Stat(ci.source); err != nil || !f.IsDir() {
			return fmt.Errorf("%s invalid", ci.source)
		}
		if err := builder.Packager.Import(ci.packageName, ci.source); err != nil {
			return err
		}
	}
	return nil
}

func (ci *CommandImport) Usage() string {
	return `gokit import 导入包
Usage:
	gokit import [options] [package]
Options:
	--source [dir]   从本地源代码导入一个包
	--insecure       从非安全连接下载包
	--help, -h       帮助信息
`
}
