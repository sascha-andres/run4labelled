package main

import (
	"flag"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"livingit.de/code/run4labelled"
	"os"
	"os/exec"
	"strings"
)

func main() {

	var configFileName = flag.String("config", "", "configuration file")
	flag.Parse()

	if *configFileName == "" {
		panic("no configuration file found")
	}

	data, err := ioutil.ReadFile(*configFileName)
	if err != nil {
		panic(err)
	}

	var cfg run4labelled.Configuration
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rec := make(chan run4labelled.Execute)
	cfg.SetChannel(rec)

	go func() {
		if err = cfg.Walk(dir); err != nil {
			panic(err)
		}
	}()

	for i := range rec {
		params := strings.Split(i.Command, " ")
		c := exec.Cmd{
			Dir:    i.Directory,
			Path:   params[0],
			Args:   params[0:],
			Stdout: os.Stdout,
		}
		err := c.Run()
		if err != nil {
			panic(err)
		}
	}
}
