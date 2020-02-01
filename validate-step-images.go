package main

import (

	// "github.com/docker/distribution/digest"
	// "github.com/docker/distribution/manifest"
	// "github.com/docker/libtrust"
	//"flag"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

	stepsFound := readJson()

	fmt.Printf("Found %d steps\n", len(stepsFound))

}

func readJson() []StepDetails {

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

	return p.finishedSteps

}

func visitMap(stepContent map[string]interface{}, p *ParsingContext) {

	for k, v := range stepContent {
		// fmt.Printf("Inside %s at %d\n", k, p.nestingLevel)

		switch v := v.(type) {
		case string:
			fmt.Println(k, v, "(Found string at )", p.nestingLevel)
			storeStepInfo(k, v, p)
		case []interface{}:
			fmt.Println(k, "(Found array):")
			storeSourcesInfo(k, v, p)
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
func storeSourcesInfo(key string, values []interface{}, p *ParsingContext) {
	if p.nestingLevel != 1 || key != "sources" || len(values) == 0 {
		return
	}
	p.currentStep.SourceURL = values[0].(string)
}

func storeStepInfo(key string, value string, p *ParsingContext) {
	switch key {
	case "image":
		storeImageInfo(value, p)
	case "version":
		if p.nestingLevel == 0 {
			goToNextStep(p)
		} else if p.nestingLevel == 1 {
			storeVersionInfo(value, p)
		}

	case "name":
		if p.nestingLevel != 1 {
			return
		}
		storeNameInfo(value, p)
	}
}

func storeImageInfo(fullDockerImage string, p *ParsingContext) {
	var dockerImage DockerImageName

	if strings.Contains(fullDockerImage, ":") {
		imageAndTag := strings.Split(fullDockerImage, ":")
		dockerImage.BaseImage = imageAndTag[0]
		dockerImage.HasTag = true
		dockerImage.Tag = imageAndTag[1]

	} else {
		dockerImage.HasTag = false
		dockerImage.BaseImage = fullDockerImage
	}
	p.currentStep.ImagesUsed = append(p.currentStep.ImagesUsed, dockerImage)
}

func storeVersionInfo(pluginVersion string, p *ParsingContext) {
	p.currentStep.Version = pluginVersion
}

func goToNextStep(p *ParsingContext) {
	if p.currentStep != nil {
		fmt.Println(">>>>>>>>>>>>>Finished with step ", p.currentStep.Name)
		p.finishedSteps = append(p.finishedSteps, *p.currentStep)
	}
	p.currentStep = new(StepDetails)
}

func storeNameInfo(pluginName string, p *ParsingContext) {
	p.currentStep.Name = pluginName
}
