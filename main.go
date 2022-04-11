package main

import (
	"flag"
	"log"
	"os"

	"sigs.k8s.io/yaml"
)

// Config is the input configuration struct for templest
type Config struct {
	OutPath        string
	TemplatePath   string
	YamlConfigFile string
}

// the YAML key used to define the overall directory structure in the input configuration
const layoutKey string = "layout"

func main() {
	config := flag.String("config", "", "path to yaml formatted stack configuration")
	templatePath := flag.String("templates", "", "path to templates dir")
	outPath := flag.String("out", "", "path to render directories and files to")
	flag.Parse()

	if *config == "" || *outPath == "" || *templatePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	c := &Config{
		TemplatePath: *templatePath,
		OutPath:      *outPath,
	}

	// open the config file and unmarshal the YAML to a go data structure
	f, err := os.ReadFile(*config)
	if err != nil {
		log.Fatalf("Error reading yaml config file %s: %v", *config, err)
	}
	var data map[string]any
	err = yaml.Unmarshal(f, &data)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML config: %v", err)
	}

	parsedLayout, err := parseLayout(data[layoutKey].(map[string]any), nil)
	if err != nil {
		log.Fatalf("Error parsing data: %v\n", err)
	}

	err = c.walkLayout(parsedLayout, "")
	if err != nil {
		log.Fatalf("Error walking parsed YAML: %v\n", err)
	}
}
