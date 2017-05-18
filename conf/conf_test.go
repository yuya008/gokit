package conf

import (
	"testing"
	"os"
	"log"
	"path"
)

const (
	testJson = "test.json"
)

var (
	pwd string
	conf *Conf
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
	if conf, err = LoadConfFile(path.Join(pwd, testJson)); err != nil {
		t.Error(err)
	}
}

func TestConf_Dump(t *testing.T) {
	t.Log(conf.Dump())
}


