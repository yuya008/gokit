package main

import (
	"git.oschina.net/yuya008/gokit/command"
	"os"
	"log"
)

func main() {
	if err := command.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
