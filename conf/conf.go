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
	Title string				`toml:"title"`
	Binary []*BinaryConf		`toml:"binary"`
	Package *PackageConf    	`toml:"package"`
	Dependent []*DependentConf	`toml:"dependent"`
}

type BinaryConf struct {
	Name string			`toml:"name"`
	Version string		`toml:"version"`
	Debug bool			`toml:"debug,omitempty"`
	BuildFlags string	`toml:"buildFlags,omitempty"`
	OutFile string		`toml:"outFile,omitempty"`
	OsArch string       `toml:"osarch,omitempty"`
	ExeName string      `toml:"exeName,omitempty"`
}

type PackageConf struct {
	Name string 		`toml:"name"`
	Version string		`toml:"version"`
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
	if c.Title == "" {
		return errors.New("`title` no set")
	}
	if len(c.Binary) > 0 {
		for i, bin := range c.Binary {
			if bin.Name == "" {
				return fmt.Errorf("%d [[Binary]] `name` no set", i)
			}
			if bin.Version == "" {
				return fmt.Errorf("%d [[Binary]] `version` no set", i)
			}
		}
	} else if c.Package.Name == "" {
		return errors.New("[[Binary]] or [Package] no set")
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