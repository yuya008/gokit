package command

import "fmt"

type CommandVersion struct {
}

func init() {
	commands["version"] = &CommandVersion{}
}

func (cv *CommandVersion) ParseArgs(args []string) (bool, error) {
	return true, nil
}

func (cv *CommandVersion) Run() error {
	fmt.Printf("%s %s\n", ProjectName, Version)
	return nil
}

func (cv *CommandVersion) Usage() string {
	return ""
}
