package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// Read the TOML file
	tomlData, err := ioutil.ReadFile("template/config.toml")
	if err != nil {
		log.Fatalf("Error reading TOML file: %s", err)
	}

	// Convert []byte to string
	tomlString := string(tomlData)

	// Print the TOML content
	fmt.Println(tomlString)
}
