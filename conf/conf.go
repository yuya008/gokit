package conf

import (
	"os"
	"io/ioutil"
	"fmt"
	"github.com/BurntSushi/toml"
	"errors"
	"bytes"
)

type Conf struct {
	Package *PackageConf    	`toml:"package"`
	Dependent []*DependentConf	`toml:"dependent"`
}

type PackageConf struct {
	Name string			`toml:"name"`
	Version string		`toml:"version"`
	Debug bool			`toml:"debug,omitempty"`
	OutFile string		`toml:"outFile,omitempty"`
	OsArch string       `toml:"osarch,omitempty"`
	ExeName string      `toml:"exeName,omitempty"`
}

type DependentConf struct {
	Source string 		`toml:"source"`
	Version string 		`toml:"version"`
	Insecure bool       `toml:"insecure,omitempty"`
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
	p := &Conf{}
	if err := toml.Unmarshal(b, p); err != nil {
		return nil, err
	}
	if err := p.verification(); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Conf) verification() error {
	if c.Package.Name == "" {
		return errors.New("[package.name] no set")
	}
	if c.Package.Version == "" {
		return errors.New("[package.version] no set")
	}
	return nil
}

func (c *Conf) String() string {
	return fmt.Sprintf("%v", c)
}

func (c *Conf) Dump() (string, error) {
	buffer := bytes.NewBufferString("")
	encoder := toml.NewEncoder(buffer)
	if err := encoder.Encode(c); err != nil {
		return "", err
	}
	return buffer.String(), nil
}