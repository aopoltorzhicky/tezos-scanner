package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Bootstrap []string `yaml:"bootstrap"`
	Attempts  struct {
		Count   int64 `yaml:"count"`
		Timeout int64 `yaml:"timeout"`
	}
	SyncedTime   int64 `yaml:"synced_time"`
	ThreadsCount int64 `yaml:"threads_count"`
}

func getConfig(filename string) (c config, err error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return c, fmt.Errorf("File %s does not exists", filename)
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}

	if err := yaml.Unmarshal(yamlFile, &c); err != nil {
		return c, err
	}
	return
}
