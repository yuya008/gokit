package main

import (
	"log"
	"github.com/BurntSushi/toml"
	"time"
	"github.com/jack/hi/a"
	"github.com/jack/hi/b"
	"github.com/jack/hi/c"
)

func main() {
	var testSimple = `
age = 250
andrew = "gallant"
kait = "brady"
now = 1987-07-05T05:45:00Z
yesOrNo = true
pi = 3.14
colors = [
	["red", "green", "blue"],
	["cyan", "magenta", "yellow", "black"],
]

[My.Cats]
plato = "cat 1"
cauchy = "cat 2"
`

	type cats struct {
		Plato  string
		Cauchy string
	}
	type simple struct {
		Age     int
		Colors  [][]string
		Pi      float64
		YesOrNo bool
		Now     time.Time
		Andrew  string
		Kait    string
		My      map[string]cats
	}
	var val simple
	_, err := toml.Decode(testSimple, &val)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("main")
	log.Println(val)
	a.PrintA()
	b.PrintB()
	c.PrintC()
}
