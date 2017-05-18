package conf

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Conf struct {
	Name string
	Version string
	BuildConfig []BuildConfig
	Dependent []DependentPackage
}

type BuildConfig struct {
	Name string
	OutFile string
	BuildFlags []string
	Debug bool
	OsArch string
	ExeName string
}

type DependentPackage struct {
	Name string
	Version string
	Insecure bool
}

func LoadConfFile(file string) (*Conf, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	conf := &Conf{}
	if err := json.Unmarshal(b, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *Conf) Dump() string {
	return fmt.Sprintf("%v", c)
}

