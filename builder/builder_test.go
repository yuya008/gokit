package builder

import (
	"testing"
	"os"
	"log"
	"path"
)

var (
	builder *Builder
	goPath string
	pwd string
)

const (
	demoDir = "demo"
)

func init() {
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}
	goPath = path.Join(pwd, "test")
}

func TestNewBuilder(t *testing.T) {
	var err error
	builder, err = NewBuilder(goPath)
	if err != nil {
		t.Error(err)
	}
}

func TestBuilder_AddProject(t *testing.T) {
	builder.AddProject(&Project{
		Name: "github.com/jack/hi",
		Source: path.Join(pwd, demoDir),
		BuildPackages: []*BuildPackage{
			&BuildPackage{
				PackageName: "github.com/jack/hi",
				OutFile: path.Join(pwd, demoDir, "hi.out"),
				BuildFlags: []string{"-gcflags", "-N -l"},
			},
			&BuildPackage{
				PackageName: "github.com/jack/hi/foo",
				OutFile: path.Join(pwd, demoDir, "foo.out"),
				BuildFlags: []string{"-gcflags", "-N -l"},
			},
			&BuildPackage{
				PackageName: "github.com/jack/hi/a",
				OutFile: path.Join(pwd, demoDir, "a.obj"),
				BuildFlags: []string{"-gcflags", "-N -l"},
			},
		},
	})
}

func TestPackager_Pull(t *testing.T) {
	if err := builder.Packager.Pull("github.com/BurntSushi/toml", false); err != nil {
		t.Error(err)
	}
	if err := builder.Packager.Pull("git.oschina.net/yuya008/testimport", false); err != nil {
		t.Error(err)
	}
}

func TestBuilder_Build(t *testing.T) {
	ch := make(chan string)

	go func() {
		defer close(ch)
		if err := builder.Build(ch); err != nil {
			t.Error(err)
		}
	}()

	for s := range ch {
		t.Log(s)
	}
}

func TestPackager_Checkout(t *testing.T) {
	if err := builder.Packager.Checkout("git.oschina.net/yuya008/testimport", "v0.0.1"); err != nil {
		t.Error(err)
	}
}

func TestBuilder_Build2(t *testing.T) {
	TestNewBuilder(t)

	builder.AddProject(&Project{
		Name: "github.com/jack/hi",
		Source: path.Join(pwd, demoDir),
		BuildPackages: []*BuildPackage{
			&BuildPackage{
				PackageName: "github.com/jack/hi",
				OutFile: path.Join(pwd, demoDir, "hi2.out"),
				BuildFlags: []string{"-gcflags", "-N -l"},
			},
		},
	})

	ch := make(chan string)

	go func() {
		defer close(ch)
		if err := builder.Build(ch); err != nil {
			t.Error(err)
		}
	}()

	for s := range ch {
		t.Log(s)
	}
}