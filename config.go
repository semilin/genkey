package main

import (
	"encoding/json"
	"fmt"
	"github.com/hjson/hjson-go"
	"io/ioutil"
)

func ReadWeights() {
	b, err := ioutil.ReadFile("weights.hjson")
	if err != nil {
		fmt.Printf("There was an issue reading the weights file.\nPlease make sure there is a 'weights.json' in this directory.")
		panic(err)
	}

	var dat map[string]interface{}

	err = hjson.Unmarshal(b, &dat)
	if err != nil {
		panic(err)
	}

	j, _ := json.Marshal(dat)
	json.Unmarshal(j, &Weight)
}
