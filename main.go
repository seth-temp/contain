package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
	"gopkg.in/lxc/go-lxc.v2"
	"os"
)

type Config struct {
	Name      string `short:"n" long:"name" default:"default" description:"Name of the container"`
	Interface string `short:"i" long:"interface" default:"0.0.0.0"`
	Args      struct {
		Command []string `required:"yes" positional-arg-name:"command"`
	} `positional-args:"yes"`
	LXCPath  string `short:"p" long:"lxcpath"`
	Template string `short:"t" long:"template" default:"lxc-alpine"`
	Distro   string
	Release  string
	Arch     string
	Debug    bool `short:"d" long:"debug" description:"Dump all debug information"`
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

/*
func init() {

}
*/

func error_exit(exit_code int, err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(exit_code)
}

func attach(c *lxc.Container, o *lxc.AttachOptions) {
	err := c.AttachShell(*o)
	if err != nil {
		error_exit(2, err)
	}
}

func create(conf *Config) *lxc.Container {
	c, err := lxc.NewContainer(conf.Name, conf.LXCPath)
	if err != nil {
		error_exit(2, err)
	}
	c.SetLogFile("/tmp" + conf.Name + ".log")
	c.SetLogLevel(lxc.TRACE)
	return c
}

func run(c *lxc.Container, conf *Config) {
	if err := c.Start(); err != nil {
		error_exit(2, err)
	}

}

func main() {
	var conf Config

	/*
	   Input validation. Don't silently fail. Print the usage instead.
	   We can assign _ to "unparsed" later, but Args nested struct in Config
	   slurps the rest of the arguments into command.
	*/
	var parser = flags.NewParser(&conf, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stderr)
		error_exit(2, err)
	}

	options := lxc.DefaultAttachOptions
	options.ClearEnv = true

	c := create(&conf)

	if conf.Debug {
		spew.Dump(c)
		spew.Dump(conf)
	}

	run(c, &conf)
	attach(c, &options)
}
