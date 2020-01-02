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

type ParsingContext struct {
	nestingLevel  int
	currentStep   *StepDetails
	finishedSteps []StepDetails
}

func main() {
	fmt.Println("Checking node image")

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
	p := new(ParsingContext)
	p.nestingLevel = 0

	for _, stepDef := range data {
		stepContent := stepDef.(map[string]interface{})

		visitMap(stepContent, p)

	}

}

func visitMap(stepContent map[string]interface{}, p *ParsingContext) {

	for k, v := range stepContent {
		fmt.Printf("Inside %s at %d\n", k, p.nestingLevel)

		switch v := v.(type) {
		case string:
			fmt.Println(k, v, "(Found string)")
			storeStepInfo(k, v, p)
		case float64:
			fmt.Println(k, v, "(Found float64)")
		case []interface{}:
			fmt.Println(k, "(Found array):")
			// for i, u := range v {
			// 	fmt.Println("    ", i, u)
			// }
		case map[string]interface{}:
			fmt.Println(k, "(Found object): ")
			p.nestingLevel++
			visitMap(v, p)
			p.nestingLevel--
		default:
			fmt.Println(k, v, "(Found unknown)")
		}

	}

}

func storeStepInfo(key string, value string, p *ParsingContext) {
	switch key {
	case "image":
		storeImageInfo(value, p)
	case "version":
		if p.nestingLevel != 0 {
			return
		}
		storeVersionInfo(value, p)
	case "sources":
		if p.nestingLevel != 1 {
			return
		}
		storeSourcesInfo(value, p)
	case "name":
		if p.nestingLevel != 1 {
			return
		}
		storeNameInfo(value, p)
	}
}

func storeImageInfo(fullDockerImage string, p *ParsingContext) {
	var dockerImage DockerImageName
	dockerImage.BaseImage = fullDockerImage
	p.currentStep.ImagesUsed = append(p.currentStep.ImagesUsed, dockerImage)
}

func storeVersionInfo(pluginVersion string, p *ParsingContext) {
	if p.currentStep != nil {
		p.finishedSteps = append(p.finishedSteps, *p.currentStep)
	}
	p.currentStep = new(StepDetails)
	p.currentStep.Version = pluginVersion
}

func storeSourcesInfo(pluginSources string, p *ParsingContext) {
	p.currentStep.SourceURL = pluginSources
}

func storeNameInfo(pluginName string, p *ParsingContext) {
	p.currentStep.Name = pluginName
}
