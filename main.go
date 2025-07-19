package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"xray-geo-exporter/config"
	"xray-geo-exporter/utils"
)

func main() {
	initConfig()

	log.Printf("XRay Geo Exporter (https://github.com/hteppl/xray-geo-exporter/)")
	log.Printf("Service started with hostname %s", config.Hostname)

	utils.StartLogMonitor()
}

func initConfig() {
	var configPath string

	flag.StringVar(&configPath, "c", "config.yaml", "Path to the configuration file")
	flag.Parse()

	if configPath == "" {
		ex, err := os.Executable()
		if err != nil {
			log.Fatalf("Error getting executable path: %v", err)
		}
		configPath = filepath.Join(filepath.Dir(ex), "config.yaml")
	}

	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
}
