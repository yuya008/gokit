package command

import (
	"fmt"
	"os"
	"gopkg.in/urfave/cli.v2"
	"errors"
	"path"
	bu "github.com/yuya008/gokit/builder"
)

type CommandImport struct {
	packageName string
	source string
	insecure bool
	cache bool
}

func init() {
	commands = append(commands, &cli.Command{
		Name: "import",
		Usage: "import a package",
		Aliases: []string{"c"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "source",
				Usage: "from a local source code import a package",
			},
			&cli.BoolFlag{
				Name: "insecure",
				Usage: "from a secure connection to download package",
			},
			&cli.BoolFlag{
				Name: "cache",
				Usage: "only import to cache",
			},
		},
		Action: func(context *cli.Context) error {
			ci := &CommandImport{}
			ci.source = context.String("source")
			ci.insecure = context.Bool("insecure")
			ci.cache = context.Bool("cache")
			args := context.Args()
			if args.Len() <= 0 {
				return errors.New("You must specify the package name")
			}
			ci.packageName = args.First()
			return ci.Run()
		},
	})
}

func (ci *CommandImport) Run() error {
	if ci.source == "" {
		fmt.Printf("pull package [%s] \n", ci.packageName)
		if err := builder.Packager.Pull(ci.packageName, ci.insecure); err != nil {
			if err != bu.PkgExist {
				return err
			}
			fmt.Printf("already cached [%s] \n", ci.packageName)
		}
	} else {
		if f, err := os.Stat(ci.source); err != nil || !f.IsDir() {
			return fmt.Errorf("%s invalid", ci.source)
		}
		fmt.Printf("import package to cache [%s] \n", ci.packageName)
		if err := builder.Packager.Import(ci.packageName, ci.source); err != nil {
			return err
		}
	}
	if !ci.cache {
		pkg, ok := builder.Packager.Lookup(ci.packageName)
		if !ok {
			return fmt.Errorf("%s not found", ci.packageName)
		}
		fmt.Printf("import project to vendor dir [%s] \n", ci.packageName)
		pkg.CopyTo(path.Join(pwd, "vendor", ci.packageName))
	}
	return nil
}
