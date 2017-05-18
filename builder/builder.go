package builder

import (
	"path/filepath"
	"os"
	"fmt"
	"errors"
)

type Builder struct {
	goPath string
	srcDir string
	projects map[string]*Project
	buildFlags []string
	Packager *Packager
}

type Project struct {
	Name string
	Source string
	linkPath string
	BuildPackages []*BuildPackage
}

type BuildPackage struct {
	PackageName string
	OutFile string
	BuildFlags []string
}

func NewBuilder(goPath string) (*Builder, error) {
	src := filepath.Join(goPath, "src")
	if err := os.MkdirAll(src, dirMode); err != nil {
		return nil, err
	}
	packager, err := NewPackager(goPath)
	if err != nil {
		return nil, err
	}
	return &Builder{
		goPath: goPath,
		srcDir: src,
		projects: make(map[string]*Project),
		Packager: packager,
	}, nil
}

func (b *Builder) AddProject(project *Project) error {
	if project.Name == "" || project.Source == "" {
		return errors.New("invalid package")
	}
	if _, ok := b.projects[project.Name]; ok {
		return fmt.Errorf("%s exist", project.Name)
	}
	if f, err := os.Stat(project.Source); err != nil || !f.IsDir() {
		return fmt.Errorf("%s not a dir", project.Source)
	}
	if err := b.linkProject(project); err != nil {
		return err
	}
	b.projects[project.Name] = project
	return nil
}

func (b *Builder) Build(c chan <- string) error {
	for _, project := range b.projects {
		if err := func() error {
			for _, pkg := range project.BuildPackages {
				if err := b.toBuild(pkg.BuildFlags, pkg.PackageName, pkg.OutFile, c); err != nil {
					return err
				}
			}
			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) Clean() {
	for _, project := range b.projects {
		os.RemoveAll(project.linkPath)
	}
}

func (b *Builder) linkProject(project *Project) error {
	newName := filepath.Join(b.srcDir, project.Name)
	os.RemoveAll(newName)
	if err := os.MkdirAll(filepath.Dir(newName), dirMode); err != nil {
		return err
	}
	if err := os.Symlink(project.Source, newName); err != nil {
		return err
	}
	project.linkPath = newName
	return nil
}

func (b *Builder) toBuild(buildFlags []string, packageName, outfile string, c chan <- string) error {
	goTools, err := NewGoTools(GoBuild, b.goPath)
	if err != nil {
		return err
	}
	defer func() {
		if c != nil {
			c <- goTools.String()
		}
	}()
	goTools.AddBuildFlags("-o", outfile)
	goTools.AddBuildFlags(buildFlags...)
	goTools.AddPackageNames(packageName)
	if s, ok := goTools.Run(); !ok {
		return errors.New(s)
	}
	return nil
}