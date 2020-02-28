package main

import (
	"flag"
	"fmt"
	"github.com/pkg/profile"
	ms "multiline-search"
	"os"
	"time"
)

func main() {
	defer profile.Start(profile.MemProfile).Stop()
	st := time.Now().Unix()
	var tmplPath, filePath string
	flag.StringVar(&tmplPath, "tmpl", "", "")
	flag.StringVar(&filePath, "in", "", "")
	flag.Parse()

	if tmplPath == "" {
		fmt.Println("template file path not specified or empty")
		return
	}
	if filePath == "" {
		fmt.Println("processing file path not specified or empty")
		return
	}

	tmpl, err := ms.ReadTemplateFromFile(tmplPath)
	if err != nil {
		fmt.Printf("cannot read template file: %s\n", err.Error())
		return
	}

	canvas, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("cannot read processing file: %s\n", err.Error())
		return
	}

	s := ms.NewSearch(canvas)
	c := s.Count(tmpl)

	fmt.Printf("Found %d bugs in %d seconds\n", c, time.Now().Unix()-st)
}
