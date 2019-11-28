package main

import (
	"github.com/heroku/docker-registry-client/registry"
	// "github.com/docker/distribution/digest"
	// "github.com/docker/distribution/manifest"
	// "github.com/docker/libtrust"
	//"flag"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "strings"
)

type StepStatus int

const (
	Ok StepStatus = iota + 1
	NotOk
	Incomplete
)

type DockerImageName struct {
	BaseImage string
	Tag       string
	HasTag    bool
}
type StepDetails struct {
	Name       string
	Version    string
	SourceURL  string
	Status     StepStatus
	ImagesUsed []DockerImageName
}

func main() {
	fmt.Println("Checking node image")

	// url      := "https://registry-1.docker.io/"
	url := "https://registry.hub.docker.com"
	username := "" // anonymous
	password := "" // anonymous
	hub, err := registry.New(url, username, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection successful")

	tags, err := hub.Tags("heroku/cedar")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("2d: ", tags)

	readJson()

}

func readJson() {

	fileName := "input.json"
	jsonFile, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var v interface{}
	json.Unmarshal(jsonData, &v)

	data := v.([]interface{})

	for _, stepDef := range data {
		stepContent := stepDef.(map[string]interface{})

		for k, v := range stepContent {
			switch v := v.(type) {
			case string:
				fmt.Println(k, v, "(Found string)")
			case float64:
				fmt.Println(k, v, "(Found float64)")
			case []interface{}:
				fmt.Println(k, "(Found array):")
				for i, u := range v {
					fmt.Println("    ", i, u)
				}
			case map[string]interface{}:
				fmt.Println(k, "(Found object): ")
				// for i, u := range v {
				// 	fmt.Println("    ", i, u)
				// }	
			default:
				fmt.Println(k, v, "(Found unknown)")
			}
		}

	}

}
