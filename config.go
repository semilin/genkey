package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func ReadWeights() {
	b, err := os.ReadFile("config.toml")
	if err != nil {
		fmt.Printf("There was an issue reading the weights file.\nPlease make sure there is a 'weights.json' in this directory.")
		panic(err)
	}

	_, err = toml.Decode(string(b), &Config)
	if err != nil {
		panic(err)
	}
}
