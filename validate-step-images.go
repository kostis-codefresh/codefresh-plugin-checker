package main

import (
	// "github.com/docker/distribution/digest"
	// "github.com/docker/distribution/manifest"
	// "github.com/docker/libtrust"
	//"flag"
	"fmt"
)

func main() {
	fmt.Println("Checking node image")

	// url := "https://registry.hub.docker.com"
	// username := "" // anonymous
	// password := "" // anonymous
	// hub, err := registry.New(url, username, password)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Connection successful")

	// tags, err := hub.Tags("heroku/cedar")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("2d: ", tags)

	stepsFound := readJSON("input.json")

	fmt.Printf("Found %d steps\n", len(stepsFound))

}
