package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
	"gopkg.in/lxc/go-lxc.v2"
	"os"
	//"strings"
)

// Config
type Config struct {
	Name      string `short:"n" long:"name" description:"Specify the name of the container"`
	Interface string `short:"i" long:"interface" default:"0.0.0.0"`

	Args struct {
		Command []string `required:"yes" positional-arg-name:"command"`
	} `positional-args:"yes"`

	LXCPath string `short:"p" long:"lxcpath" description:"Specify container path"`
	// Alpine is all the container OS rage these days
	Template string `short:"t" long:"template" default:"/usr/share/lxc/templates/lxc-alpine"`

	/*
		We probably don't need all this
		Distro     string `short:"d" long:"distro" default:"alpine" description:"Distro for the template"`
		Release    string `short:"r" long:"release" default:"v3.3" description:"Release for the template"`
		Arch       string `short:"a" long:"arch" default:"amd64" description:"Arch for the template"`
		FlushCache bool `short:"C" long:"flush-cache" description:"Flush LXC cache for image"`
		Validation bool `short:"V" long:"validation" description:"GPG Validation"`
	*/
	Interactive bool `short:"I" long:"interactive" description:"Attach TTY"`
	Debug       bool `short:"D" long:"debug" description:"Dump all debug information"`
	Help        bool `short:"h" long:"help" description:"Show this help message"`
}

func errorExit(exit_code int, err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(exit_code)
}

func attach(c *lxc.Container, o *lxc.AttachOptions) {
	err := c.AttachShell(*o)
	if err != nil {
		errorExit(2, err)
	}
}

func create(conf *Config) *lxc.Container {
	c, err := lxc.NewContainer(conf.Name, conf.LXCPath)
	if err != nil {
		fmt.Printf("FOOOOO")
		errorExit(2, err)
	}
	c.SetLogFile("/tmp" + conf.Name + ".log")
	c.SetLogLevel(lxc.TRACE)
	options := lxc.TemplateOptions{
		Template: conf.Template,
	}
	//if !(c.Defined()) {
	if err := c.Create(options); err != nil {
		fmt.Printf("Could not create container \"%s\"\n", conf.Name)
		errorExit(2, err)
	}
	//	}
	return c
}

/*
// Whoops, might have been a little confused with the verbage
// run !+ execute
func run(c *lxc.Container, conf *Config) {
	cmd := strings.Join(conf.Args.Command, " ")
	fmt.Printf("Starting container \"%s\"...\n", conf.Name)
	if err := c.Start(); err != nil {
		fmt.Printf("Failed to run container with command \"%s\"\n", cmd)
		errorExit(2, err)
	}
}
*/

func exec(c *lxc.Container, conf *Config) {
	if output, err := c.Execute(conf.Args.Command...); err != nil {
		errorExit(2, err)
	} else {
		fmt.Printf("%s", output)
	}
}

func validateConfig(conf *Config) {
	if conf.LXCPath == "" {
		conf.LXCPath = lxc.DefaultConfigPath()
	}

	if conf.Name == "" {
		conf.Name = randomdata.SillyName()
	}
}

func checkTemplateExistence(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Could not stat LXC template \"%s\"\n", path)
		fmt.Printf("Ensure lxc packages are installed on your system\n")
		errorExit(2, err)
	}
}

func printHelp(parser *flags.Parser) {
	//fmt.Printf("%s\n", unparsed)
	parser.WriteHelp(os.Stderr)
	os.Exit(0)
}

func main() {
	var conf Config

	/*
	   Input validation. Don't silently fail. Print the usage instead.
	   We can assign _ to "unparsed" later, but Args nested struct in Config
	   slurps the rest of the arguments into command.
	*/

	var parser = flags.NewParser(&conf, flags.Default)
	unparsed, err := parser.Parse()
	if err != nil || len(unparsed) > 1 || conf.Help {
		printHelp(parser)
		//errorExit(2, err)
	}

	validateConfig(&conf)

	checkTemplateExistence(conf.Template)

	options := lxc.DefaultAttachOptions
	options.ClearEnv = true

	c := create(&conf)

	if conf.Debug {
		spew.Dump(c)
		spew.Dump(conf)
	}

	//run(c, &conf)
	exec(c, &conf)

	if conf.Interactive {
		attach(c, &options)

	}
}
