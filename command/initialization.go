package command

import (
	"gopkg.in/urfave/cli.v2"
	"bufio"
	"os"
	"fmt"
	"strings"
	"github.com/yuya008/gokit/conf"
	"github.com/yuya008/gokit/dep"
	"runtime"
	"path"
	"io/ioutil"
)

type Initialization struct {
	name string
	version string
}

func init() {
	commands = append(commands, &cli.Command{
		Name: "init",
		Usage: "initialized to gokit project",
		Aliases: []string{"i"},
		Action: func(context *cli.Context) error {
			init := &Initialization{
				name: context.String("name"),
				version: context.String("version"),
			}
			return init.Run()
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "name",
				Usage: "set the package name(e.g. github.com/foo/bar)",
			},
			&cli.StringFlag{
				Name: "version",
				Usage: "set the package version(e.g. v1.0.0)",
			},
		},
	})
}

func (init *Initialization) Run() error {
	init.askInitInfo()
	c := &conf.Conf{
		Package: &conf.PackageConf{
			Name: init.name,
			Version: init.version,
			Debug: true,
		},
	}
	depPackage, err := scanDependencyPackage(init.name)
	if err != nil {
		return err
	}
	fmt.Printf("analysis dependency [%s]\n", init.name)
	for _, dp := range depPackage {
		fmt.Println(`  ->`, dp)
		c.Dependent = append(c.Dependent, &conf.DependentConf{
			Source: dp,
		})
	}
	s, err := c.Dump()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(pwd, confFileName), []byte(s), 0766)
}

func (init *Initialization) askInitInfo() {
	if init.name == "" {
		fmt.Println("Please enter a package name(e.g. github.com/foo/bar):")
		readAnswer(&init.name)
	}
	if init.version == "" {
		fmt.Println("Please enter a package version(e.g. v1.0.0):")
		readAnswer(&init.version)
	}
}

func scanDependencyPackage(name string) ([]string, error) {
	scanner := dep.DepScanner{
		RootDir: pwd,
		Deep: true,
		HasRoot: true,
	}
	packages, err := scanner.Scan()
	if err != nil {
		return nil, err
	}
	var pkg map[string]bool = make(map[string]bool)
	for _, p := range packages {
		for _, imp := range p.Imports {
			pkg[imp] = true
		}
		for _, timp := range p.TestImports {
			pkg[timp] = true
		}
	}
	var ret []string
	for k := range pkg {
		if isStdLib(k) || isInnerPkg(name, k) {
			continue
		}
		ret = append(ret, k)
	}
	return ret, nil
}

func readAnswer(s *string) {
	var err error
	reader := bufio.NewReader(os.Stdin)
	for {
		if *s, err = reader.ReadString('\n'); err != nil {
			fmt.Println(err)
			continue
		}
		*s = strings.TrimSpace(*s)
		if *s == "" {
			continue
		}
		return
	}
}

func isStdLib(pkgName string) bool {
	goRoot := runtime.GOROOT()
	if f, err := os.Stat(path.Join(goRoot, "src", pkgName)); err == nil && f.IsDir() {
		return true
	}
	return false
}

func isInnerPkg(name, pkgName string) bool {
	return strings.HasPrefix(pkgName, name)
}