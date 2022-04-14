package main

import (
	"flag"
	"log"
	"os"

	templest "github.com/rorski/templest/internal"
)

func main() {
	config := flag.String("config", "", "path to yaml formatted stack configuration")
	templatePath := flag.String("templates", "", "path to templates dir")
	outPath := flag.String("out", "", "path to render directories and files to")
	flag.Parse()

	if *config == "" || *outPath == "" || *templatePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	c := &templest.Config{
		YAMLConfigFile: *config,
		TemplatePath:   *templatePath,
		OutPath:        *outPath,
	}

	err := templest.Run(c)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
