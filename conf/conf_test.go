package conf

import (
	"testing"
	"os"
	"log"
	"path"
)

const (
	testJson = "test.toml"
)

var (
	pwd string
	c *Conf
)

func init() {
	var err error
	if pwd, err = os.Getwd(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func TestLoadConfFile(t *testing.T) {
	var err error
	if c, err = LoadConfFile(path.Join(pwd, testJson)); err != nil {
		t.Error(err)
	} else {
		t.Log(c)
	}
}

func TestConf_Dump(t *testing.T) {
	c := Conf{}
	s, err := c.Dump()
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}

func TestNewBinaryConfTemplate(t *testing.T) {
	c := NewBinaryConfTemplate()
	s, err := c.Dump()
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}

func TestNewPackageConfTemplate(t *testing.T) {
	c := NewPackageConfTemplate()
	s, err := c.Dump()
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}
