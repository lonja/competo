package main

import (
	"flag"
	"fmt"
	ms "github.com/lonja/multiline-search"
	"os"
)

type CountCommand struct {
	fs       *flag.FlagSet
	tmplName string
	fileName string
}

func NewCountCommand() *CountCommand {
	fs := flag.NewFlagSet("count", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Print(`Usage of count:
  count -tmpl <template_file> -in <file_to_process>
Flags:
`)
		fs.PrintDefaults()
		os.Exit(2)
	}

	cc := &CountCommand{
		fs: fs,
	}

	fs.StringVar(&cc.tmplName, "tmpl", "./tmpl.txt", "Template file search templates in")
	fs.StringVar(&cc.fileName, "file", "./landscape.txt", "File to find and count templates")

	return cc
}

func (c *CountCommand) Run() error {
	tmpl, err := ms.ReadTemplateFromFile(c.tmplName)
	if err != nil {
		return fmt.Errorf("cannot read template file: %s", err.Error())
	}

	f, err := os.Open(c.fileName)
	if err != nil {
		return fmt.Errorf("cannot read source file: %ss", err.Error())
	}

	count := ms.NewSearch(f).Count(tmpl)

	fmt.Printf("Found %d bugs\n", count)

	return nil
}

func (c *CountCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *CountCommand) Name() string {
	return c.fs.Name()
}

func (c *CountCommand) Usage() {
	c.fs.Usage()
}

type Runner interface {
	Run() error
	Init(args []string) error
	Name() string
	Usage()
}

func main() {
	cmds := []Runner{
		NewCountCommand(),
	}

	if len(os.Args) <= 1 {
		fmt.Print(`usage: multiline-search <command> [<args>] 
	count	 Count all occurences of template in file`)
		os.Exit(2)
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			if err := cmd.Init(os.Args[2:]); err != nil {
				fmt.Printf("error parsing arguments: %s\n", err.Error())
				cmd.Usage()
				return
			}
			if err := cmd.Run(); err != nil {
				fmt.Println(err.Error())
				cmd.Usage()
			}
		}
	}
}
