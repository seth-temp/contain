package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
	"gopkg.in/lxc/go-lxc.v2"
	"os"
	//"strings"
)

// Config
type Config struct {
	Name      string `short:"n" long:"name" default:"default" description:"Name of the container"`
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
}

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
		fmt.Printf("FOOOOO")
		error_exit(2, err)
	}
	c.SetLogFile("/tmp" + conf.Name + ".log")
	c.SetLogLevel(lxc.TRACE)
	options := lxc.TemplateOptions{
		Template: conf.Template,
	}
	if !(c.Defined()) {
		if err := c.Create(options); err != nil {
			fmt.Printf("Could not create container \"%s\"\n", conf.Name)
			error_exit(2, err)
		}
	}
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
		error_exit(2, err)
	}
}
*/

func exec(c *lxc.Container, conf *Config) {
	if output, err := c.Execute(conf.Args.Command...); err != nil {
		error_exit(2, err)
	} else {
		fmt.Printf("%s", output)
	}

}

func validate_config(conf *Config) {
	if conf.LXCPath == "" {
		conf.LXCPath = lxc.DefaultConfigPath()
	}
}

func check_template_existence(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Could not stat LXC template \"%s\"\n", path)
		fmt.Printf("Ensure lxc packages are installed on your system\n")
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

	validate_config(&conf)

	check_template_existence(conf.Template)

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
