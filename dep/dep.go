package dep

import (
	"go/build"
	"path/filepath"
	"os"
	"errors"
)

type DepScanner struct {
	RootDir string
	Deep bool
	HasRoot bool
}

var BreakErr = errors.New("break!")

func (ds *DepScanner) Scan() ([]*build.Package, error) {
	var packages []*build.Package
	var err error
	if ds.RootDir, err = filepath.Abs(ds.RootDir); err != nil {
		return nil, err
	}
	if f, err := os.Stat(ds.RootDir); err != nil || !f.IsDir() {
		return nil, errors.New("RootDir invalid")
	}
	if ds.Deep {
		if packages, err = ds.deepScan(); err != nil {
			return nil, err
		}
	} else {
		p, err := scan(ds.RootDir)
		if err != nil {
			return nil, err
		}
		packages = append(packages, p)
	}
	return packages, nil
}

func scan(dir string) (*build.Package, error) {
	p, err := build.ImportDir(dir, 0)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (ds *DepScanner) deepScan() ([]*build.Package, error) {
	var packages []*build.Package
	err := filepath.Walk(ds.RootDir, func (path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		if path == ds.RootDir && !ds.HasRoot {
			return nil
		}
		if p, err := scan(path); err == nil {
			packages = append(packages, p)
		} else {
			return filepath.SkipDir
		}
		return nil
	})
	if err == BreakErr {
		return packages, nil
	}
	if err != nil {
		return nil, err
	}
	return packages, nil
}
