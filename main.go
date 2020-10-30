package main

import (
	"log"
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"github.com/melbahja/bundler/bundle"
)


var (
	cfg = flag.String("config", "bundler.yaml", "Bundler config file.")
)

func init() {
	flag.Parse()
}

func main() {

	data, err := ioutil.ReadFile(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	bundler := &bundle.Bundler{}

	if err = yaml.Unmarshal(data, bundler); err != nil {
		log.Fatal(err)
	}

	if err = bundler.Run(); err != nil {
		log.Fatal(err)
	}
}
