package builder

import (
	"os"
	"github.com/tadvi/rkv"
	"path"
	"fmt"
	"os/exec"
	"errors"
	"path/filepath"
	"strings"
	"io"
)

type Packager struct {
	dir string
	kvStorage *rkv.Rkv
}

const (
	dirMode = 0755
	dataFileName = ".db"
)

func NewPackager(dir string) (*Packager, error) {
	if err := os.MkdirAll(dir, dirMode); err != nil {
		return nil, err
	}
	kv, err := rkv.New(path.Join(dir, dataFileName))
	if err != nil {
		return nil, err
	}
	return &Packager{
		dir: dir,
		kvStorage: kv,
	}, nil
}

func (p *Packager) Import(packageName, source string) error {
	if packageName == "" || source == "" {
		return nil
	}
	if err := checkFromLocalSourcePackage(source); err != nil {
		return err
	}
	if p.kvStorage.Exist(packageName) {
		return fmt.Errorf("%s already exist", packageName)
	}
	packagePath := path.Join(p.dir, packageName)
	if err := dirCopy(packagePath, source); err != nil {
		return err
	}
	pkg := &Package{
		PackageName: packageName,
		PackageSource: packagePath,
	}
	p.kvStorage.Put(packageName, pkg)
	return nil
}

func (p *Packager) Pull(packageName string, insecure bool) error {
	if packageName == "" {
		return nil
	}
	if p.kvStorage.Exist(packageName) {
		return fmt.Errorf("%s already exist", packageName)
	}
	tools, err := NewGoTools(GoGet, p.dir)
	if err != nil {
		return err
	}
	tools.AddBuildFlags("-d", "-u")
	if insecure {
		tools.AddBuildFlags("-insecure")
	}
	tools.AddPackageNames(packageName)
	if s, ok := tools.Run(); !ok {
		return errors.New(s)
	}
	pkg := &Package{
		PackageName: packageName,
		PackageSource: path.Join(p.dir, "src", packageName),
	}
	p.kvStorage.Put(packageName, pkg)
	return nil
}

func (p *Packager) Checkout(packageName string, v string) error {
	if packageName == "" || v == "" {
		return nil
	}
	if !p.kvStorage.Exist(packageName) {
		return fmt.Errorf("%s not exist", packageName)
	}
	pkg := &Package{}
	if err := p.kvStorage.Get(packageName, pkg); err != nil {
		return err
	}
	if f, err := os.Stat(pkg.PackageSource); err != nil || !f.IsDir() {
		p.kvStorage.Delete(packageName)
		return fmt.Errorf("%s not found", packageName)
	}
	if s, err := gitCheckOut(pkg.PackageSource, v); err != nil {
		return errors.New(s)
	}
	return nil
}

func (p *Packager) Lookup(packageName string) (*Package, bool) {
	pkg := &Package{}
	if err := p.kvStorage.Get(packageName, pkg); err != nil {
		return nil, false
	}
	return pkg, true
}

type Package struct {
	PackageName string
	PackageSource string
}

func checkFromLocalSourcePackage(source string) error {
	if f, err := os.Stat(path.Join(source)); err != nil || !f.IsDir() {
		return err
	}
	if f, err := os.Stat(path.Join(source, ".git")); err != nil || !f.IsDir(){
		return fmt.Errorf("%s invalid git source", source)
	}
	if _, err := gitCheckOut(source, "master"); err != nil {
		return errors.New("source can't git checkout master")
	}
	return nil
}

func gitCheckOut(rootDir, v string) (string, error) {
	git, err := FindGitTool()
	if err != nil {
		return "", err
	}
	cmd := exec.Command(git, "checkout", v)
	cmd.Dir = rootDir
	s, err := cmd.CombinedOutput()
	return string(s), err
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