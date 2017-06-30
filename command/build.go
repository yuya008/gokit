package command

import (
	"fmt"
	"github.com/yuya008/gokit/conf"
	//"path"
	"strings"
	bu "github.com/yuya008/gokit/builder"
	//"os"
	"gopkg.in/urfave/cli.v2"
	"path"
	"os"
)

type CommandBuild struct {
	confFile string
	conf *conf.Conf
	release bool
}

const (
	confFileName = "gokit.toml"
	targetDir = "_target"
)

func init() {
	commands = append(commands, &cli.Command{
		Name: "build",
		Usage: "build project",
		Aliases: []string{"b"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "conf",
				Usage: "specify the configuration file",
			},
			&cli.BoolFlag{
				Name: "release",
				Usage: "specify the compilation mode for release",
			},
		},
		Action: func(context *cli.Context) error {
			cb := &CommandBuild{}
			cb.confFile = context.String("conf")
			cb.release = context.Bool("release")
			return cb.Run()
		},
	})
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
	if len(cb.conf.Binary) <= 0 {
		return nil
	}
	if cb.release {
		for i := 0; i < len(cb.conf.Binary); i++ {
			cb.conf.Binary[i].Debug = false
		}
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

func loadAndCheckConffile(file string) (*conf.Conf, error) {
	if file == "" {
		file = path.Join(pwd, confFileName)
	}
	c, err := conf.LoadConfFile(file)
	if err != nil {
		return nil, newConfError(err.Error())
	}
	return c, nil
}

func newConfError(format string, a ...interface{}) error {
	return fmt.Errorf(confFileName + " : " + format, a...)
}

func createBuildPackage(buildConfig *conf.BinaryConf) *bu.BuildPackage {
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
		OsArch: buildConfig.OsArch,
	}
	if bp.OutFile == "" {
		if buildConfig.ExeName != "" {
			bp.OutFile = path.Join(pwd, targetDir, buildConfig.Version, mode, buildConfig.ExeName)
		} else {
			packageNameSlice := packageNameSplit(bp.PackageName)
			bp.OutFile = path.Join(pwd, targetDir, buildConfig.Version, mode, packageNameSlice[len(packageNameSlice) - 1])
		}
		os.MkdirAll(path.Dir(bp.OutFile), dirMode)
	}
	return bp
}

func createProject(conf *conf.Conf, src string) *bu.Project {
	project := &bu.Project{}
	project.Name = conf.Binary[0].Name
	project.Source = src

	for _, bc := range conf.Binary {
		project.BuildPackages = append(
			project.BuildPackages,
			createBuildPackage(bc),
		)
	}
	return project
}

func dependentHandle(conf *conf.Conf) error {
	for _, dep := range conf.Dependent {
		if _, ok := builder.Packager.Lookup(dep.Source); !ok {
			fmt.Printf("download --> %s\n", dep.Source)
			if err := builder.Packager.Pull(dep.Source, dep.Insecure); err != nil {
				fmt.Println(err)
			}
		}
	}
	for _, dep := range conf.Dependent {
		if dep.Version != "" {
			fmt.Printf("process ---> %s [%s]\n", dep.Source, dep.Version)
			if pkg, ok := builder.Packager.Lookup(dep.Source); ok {
				destDir := path.Join(pwd, "vendor", dep.Source)
				if _, err := os.Stat(destDir); err == nil {
					os.RemoveAll(destDir)
				}
				if err := pkg.Checkout(dep.Version); err != nil {
					return fmt.Errorf("checkout %s failure", dep.Version)
				}
				if err := pkg.CopyTo(destDir); err != nil {
					return err
				}
				pkg.Checkout("master")
			} else {
				fmt.Printf("%s not found\n", dep.Source)
			}
		} else {
			fmt.Printf("process ---> %s\n", dep.Source)
		}
	}
	return nil
}

func packageNameSplit(s string) []string {
	return strings.Split(s, "/")
}