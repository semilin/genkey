package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func findConfig() string {
	if fileExists("./config.toml") {
		return "config.toml"
	}
	config_dir := os.Getenv("XDG_CONFIG_HOME")
	path := filepath.Join(config_dir, "genkey", "config.toml")
	if fileExists(path) {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	path = filepath.Join(home, ".config", "genkey", "config.toml")
	if fileExists(path) {
		return path
	}

	println("Couldn't find config.toml in any of local directory, $XDG_CONFIG_HOME/genkey/config.toml, or ~/.config/genkey/config.toml.")
	os.Exit(1)
	return ""
}

func ReadWeights() {
	path := findConfig()
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("There was an issue reading the config file.")
		panic(err)
	}

	_, err = toml.Decode(string(b), &Config)
	if err != nil {
		panic(err)
	}

	if !fileExists(filepath.Join(Config.Paths.Corpora, Config.Corpus) + ".json") {
		fmt.Printf("Invalid config: Corpus [%s] does not exist.\n", Config.Corpus)
		os.Exit(1)
	}

	if Config.Generation.Selection > Config.Generation.InitialPopulation {
		fmt.Println("Invalid config: Generation.Selection cannot be greater than Generation.InitialPopulation.")
		os.Exit(1)
	}
}
