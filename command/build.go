package command

import (
	"fmt"
	"git.oschina.net/yuya008/gokit/conf"
	"path"
	"strings"
	bu "git.oschina.net/yuya008/gokit/builder"
	"os"
	"path/filepath"
	"io"
)

type CommandBuild struct {
	confFile string
	conf *conf.Conf
}

var platforms = map[string]bool {
	"darwin/386": true,
	"darwin/amd64": true,
	"linux/386": true,
	"linux/amd64": true,
	"linux/arm": true,
	"freebsd/386": true,
	"freebsd/amd64": true,
	"openbsd/386": true,
	"openbsd/amd64": true,
	"windows/386": true,
	"windows/amd64": true,
	"freebsd/arm": true,
	"netbsd/386": true,
	"netbsd/amd64": true,
	"netbsd/arm": true,
	"plan9/386": true,
}

const (
	confFileName = "gokit-conf.json"
	targetDir = "target"
)

func (cb *CommandBuild) ParseArgs(args []string) (ok bool, error error) {
	var option string
	defer func() {
		if err := recover(); err != nil {
			ok = false
			error = fmt.Errorf("'%s' option value invalid", option)
		}
	}()
	for i := 0; i < len(args); i++ {
		option = args[i]
		switch option {
		case "-h", "--help":
			return false, nil
		case "--conf":
			i++
			cb.confFile = args[i]
		}
	}
	return true, nil
}

func (cb *CommandBuild) Run() error {
	var err error
	cb.conf, err = loadAndCheckConffile(cb.confFile)
	if err != nil {
		return err
	}
	if err := dependentHandle(cb.conf); err != nil {
		return err
	}
	if err := builder.AddProject(createProject(cb.conf, pwd)); err != nil {
		return err
	}
	ch := make(chan string)
	go func() {
		defer close(ch)
		if err := builder.Build(ch); err != nil {
			ch <- err.Error()
		}
	}()
	for s := range ch {
		fmt.Println(s)
	}
	builder.Clean()
	return nil
}

func (cb *CommandBuild) Usage() string {
	return `gokit build 编译构建包
Usage:
	gokit build [options]
Option:
	--help, -h    帮助信息
	--conf [file] 构建配置文件
`
}

func loadAndCheckConffile(file string) (*conf.Conf, error) {
	if file == "" {
		file = path.Join(pwd, confFileName)
	}
	fmt.Printf("加载配置文件 %s\n", file)
	c, err := conf.LoadConfFile(file)
	if err != nil {
		return nil, newConfError(err.Error())
	}
	if c.Name == "" {
		return nil, newConfError("'%s' is not defined", "name")
	}
	if c.Version == "" {
		return nil, newConfError("'%s' is not defined", "version")
	}
	if len(c.BuildConfig) == 0 {
		return nil, newConfError("'%s' is not defined", "buildconfig")
	}
	for index, config := range c.BuildConfig {
		if config.Name == "" {
			return nil, newConfError("'%s[%d]' 'name' is not defined", "buildconfig", index)
		}
		if config.OsArch != "" {
			config.OsArch = strings.ToLower(config.OsArch)
			if _, ok := platforms[config.OsArch]; !ok {
				return nil, newConfError("'%s' invalid", config.OsArch)
			}
		}
	}
	return c, nil
}

func newConfError(format string, a ...interface{}) error {
	return fmt.Errorf(confFileName + " : " + format, a...)
}

func createBuildPackage(buildConfig conf.BuildConfig, version string) *bu.BuildPackage {
	var buildFlags []string
	var mode string
	if buildConfig.Debug {
		buildFlags = []string{"-gcflags", "-N -l"}
		mode = "debug"
	} else {
		mode = "release"
	}
	bp := &bu.BuildPackage{
		PackageName: buildConfig.Name,
		BuildFlags: buildFlags,
		OutFile: buildConfig.OutFile,
	}
	if bp.OutFile == "" {
		bp.OutFile = path.Join(pwd, targetDir, version, mode, path.Base(bp.PackageName))
		os.MkdirAll(path.Dir(bp.OutFile), dirMode)
	}
	return bp
}

func createProject(conf *conf.Conf, src string) *bu.Project {
	project := &bu.Project{}
	project.Name = conf.Name
	project.Source = src

	for _, bc := range conf.BuildConfig {
		project.BuildPackages = append(
			project.BuildPackages,
			createBuildPackage(bc, conf.Version),
		)
	}
	return project
}

func dependentHandle(conf *conf.Conf) error {
	for _, dep := range conf.Dependent {
		if _, ok := builder.Packager.Lookup(dep.Name); !ok {
			fmt.Printf("download --> %s\n", dep.Name)
			if err := builder.Packager.Pull(dep.Name, dep.Insecure); err != nil {
				fmt.Println(err)
			}
		}
	}
	if len(conf.Dependent) > 0 {
		fmt.Println()
	}
	for _, dep := range conf.Dependent {
		fmt.Printf("---> %s\n", dep.Name)
		if dep.Version != "" {
			if pkg, ok := builder.Packager.Lookup(dep.Name); ok {
				destDir := path.Join(pwd, "vendor", dep.Name)
				if _, err := os.Stat(destDir); err == nil {
					os.RemoveAll(destDir)
				}
				if err := builder.Packager.Checkout(dep.Name, dep.Version); err != nil {
					return fmt.Errorf("checkout %s failure", dep.Version)
				}
				if err := dirCopy(destDir, pkg.PackageSource); err != nil {
					return err
				}
				builder.Packager.Checkout(dep.Name, "master")
			} else {
				fmt.Printf("%s not found\n", dep.Name)
			}
		}
	}
	return nil
}

func dirCopy(dest, src string) error {
	os.MkdirAll(dest, dirMode)
	return filepath.Walk(src, func(file string, info os.FileInfo, err error) error {
		if src == file {
			return nil
		}
		suffix := strings.Replace(file, src, "", -1)
		newFile := path.Join(dest, suffix)
		if info.IsDir() {
			if err := os.MkdirAll(newFile, info.Mode()); err != nil {
				return err
			}
		} else {
			nf, err := os.OpenFile(newFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				return err
			}
			defer nf.Close()
			of, err := os.Open(file)
			if err != nil {
				return err
			}
			defer of.Close()
			if _, err := io.Copy(nf, of); err != nil {
				return err
			}
		}
		return nil
	})
}