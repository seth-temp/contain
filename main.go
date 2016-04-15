package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
	"gopkg.in/lxc/go-lxc.v2"
)

type Config struct {
	Name      string `short:"n" long:"name" default:"default" description:"Name of the container"`
	Interface string `short:"i" long:"interface" default:"0.0.0.0"`
	Args      struct {
		Command []string `required:"yes" positional-arg-name:"command"`
	} `positional-args:"yes"`
	LXCPath  string `short:"p" long:"lxcpath"`
	Template string ``
	Distro   string
	Release  string
	Arch     string
}

func parseArgs(c *Config) {
	/*
		c.Name = flag.String("name", "default", "The name of the container")
		c.Command = flag.Arg(0)
		flag.Parse()
		g := flag.Arg(0)
		fmt.Printf("%v\n", g)
	*/
}

func main() {
	var conf Config
	unparsed, err := flags.Parse(&conf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%v\n", unparsed)
	}
	//parseArgs(&conf)
	c, err := lxc.NewContainer(conf.Name, "/var/")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		spew.Dump(c)
	}

	//options := lxc.TemplateOptions{}
	spew.Dump(conf)

}
