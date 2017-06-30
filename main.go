package main

import (
	"os"
	"github.com/yuya008/gokit/command"
	"fmt"
)

func main() {
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
