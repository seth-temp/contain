package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
	"github.com/sethdmoore/go-lxc"
	"strings"
	//"gopkg.in/lxc/go-lxc.v2"  // bugged library, until #59 is merged. Use our fork
	"os"
)

// Config is set by command line flags mostly
type Config struct {
	Name      string `short:"n" long:"name" description:"Specify the name of the container"`
	Interface string `short:"i" long:"interface" default:"0.0.0.0"`

	Args struct {
		Command []string `required:"yes" positional-arg-name:"command"`
	} `positional-args:"yes"`

	LXCPath string `short:"p" long:"lxcpath" description:"Specify container path"`
	// Alpine is all the container OS rage these days
	Template string `short:"t" long:"template" default:"/usr/share/lxc/templates/lxc-alpine"`
	// Currently broken, not in scope to fix IMO
	Interactive bool    `short:"I" long:"interactive" description:"Attach TTY"`
	MemoryLimit float64 `short:"m" long:"memory-limit" default:"0" description:"Memory limit, in bytes"`
	Debug       bool    `short:"D" long:"debug" description:"Dump all debug information"`
	Help        bool    `short:"h" long:"help" description:"Show this help message"`
}

// errorExit exits with exit_code and prints the failure message
func errorExit(exit_code int, err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(exit_code)
}

// create provisions an lxc container (if necessary)
func create(conf *Config) *lxc.Container {
	var c *lxc.Container
	var err error

	// ensure we're not attempting to recreate the same container
	activeContainers := lxc.DefinedContainers(conf.LXCPath)
	for idx := range activeContainers {
		if activeContainers[idx].Name() == conf.Name {
			fmt.Printf("Found existing container \"%s\"\n", conf.Name)
			c = &activeContainers[idx]
		}

	}

	// If we did not find a container, create a struct for one
	if c == nil {
		c, err = lxc.NewContainer(conf.Name, conf.LXCPath)
		if err != nil {
			errorExit(2, err)
		}
	}

	// double check on whether the container is defined
	if !(c.Defined()) {
		fmt.Printf("Creating new container: %s\n", conf.Name)
		options := lxc.TemplateOptions{
			Template: conf.Template,
		}
		// provision the container
		if err = c.Create(options); err != nil {
			fmt.Printf("Could not create container \"%s\"\n", conf.Name)
			errorExit(2, err)
		}
	}

	// trace level logs end up here
	c.SetLogFile("/tmp/" + conf.Name + ".log")
	c.SetLogLevel(lxc.TRACE)

	return c
}

// exec executes a command in the context of the container
func exec(c *lxc.Container, conf *Config) {
	var output []byte
	var err error
	// stdout and stderr are unfornutately concatenated
	if output, err = c.Execute(conf.Args.Command...); err != nil {
		if len(output) != 0 {
			fmt.Printf("%s\n", output)
		}
		errorExit(2, err)
	} else {
		fmt.Printf("%s", output)
	}
}

// run starts a container
func run(c *lxc.Container, conf *Config) {
	cmd := strings.Join(conf.Args.Command, " ")
	fmt.Printf("Starting container \"%s\"...\n", conf.Name)
	if err := c.Start(); err != nil {
		fmt.Printf("Failed to run container with command \"%s\"\n", cmd)
		errorExit(2, err)
	}
}

// attach binds the TTY to the running container
func attach(c *lxc.Container, o *lxc.AttachOptions) {
	err := c.AttachShell(*o)
	if err != nil {
		errorExit(2, err)
	}
}

// parseArgs operates on a reference to Config
func parseArgs(conf *Config) {
	var parser = flags.NewParser(conf, flags.Default)

	/*
	   Input validation. Don't silently fail. Print the usage instead.
	   We might do something with "unparsed" later, but the Args nested
	   struct in Config slurps the rest of the arguments into command.

	   There seems to be a bug where --help prints twice... I tried to
	   mitigate it by overriding with my own --help. I think it's caused
	   by the fact that I have required args? Not worth investing any more
	   time
	*/

	unparsed, err := parser.Parse()
	if err != nil || len(unparsed) > 1 || conf.Help {
		printHelp(parser)
	}
}

// validateConfig sets some defaults if the flags are unset
func validateConfig(conf *Config) {
	// Hopefully lxc package derives this correctly
	if conf.LXCPath == "" {
		conf.LXCPath = lxc.DefaultConfigPath()
	}

	// Generate "Docker-style" container names if it is not provided
	if conf.Name == "" {
		conf.Name = randomdata.SillyName()
	}
}

// checkTemplateExistence simply stats a template specified
func checkTemplateExistence(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Could not stat LXC template \"%s\"\n", path)
		fmt.Printf("Ensure lxc packages are installed on your system\n")
		errorExit(2, err)
	}
}

// printHelp does what it says
func printHelp(parser *flags.Parser) {
	parser.WriteHelp(os.Stderr)
	os.Exit(0)
}

func main() {
	var conf Config

	parseArgs(&conf)

	validateConfig(&conf)

	checkTemplateExistence(conf.Template)

	options := lxc.DefaultAttachOptions
	options.ClearEnv = true

	c := create(&conf)

	if conf.Debug {
		spew.Dump(c)
		spew.Dump(conf)
	}

	if conf.MemoryLimit > 0 {
		var b lxc.ByteSize
		b = lxc.ByteSize(conf.MemoryLimit)
		c.SetMemoryLimit(b)
	}

	if conf.Interactive {
		run(c, &conf)
		attach(c, &options)

	} else {
		exec(c, &conf)
	}
}
