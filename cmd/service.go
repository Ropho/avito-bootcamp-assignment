package main

import (
	"flag"
	"log"

	"github.com/Ropho/avito-bootcamp-assignment/internal/boot"
)

const configPath = "./config"

var (
	configName string
)

func init() {
	flag.StringVar(&configName, "config_name", "config", "name for config file in folder ./config")
}

func main() {
	flag.Parse()

	log.Fatalf("server stopped: %v", boot.App(configPath, configName))
}
