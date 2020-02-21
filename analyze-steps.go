package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type stepStatus int

const (
	ok stepStatus = iota + 1
	notOk
	incomplete
)

type dockerImageName struct {
	baseImage string
	tag       string
	hasTag    bool
}
type stepDetails struct {
	name       string
	version    string
	sourceURL  string
	status     stepStatus
	imagesUsed []dockerImageName
}

type parsingContext struct {
	nestingLevel  int
	currentStep   *stepDetails
	finishedSteps []stepDetails
}

func readJSON(fileName string) []stepDetails {

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
	p := new(parsingContext)
	p.nestingLevel = 0

	for _, stepDef := range data {
		stepContent := stepDef.(map[string]interface{})

		readSingleStep(stepContent, p)

	}

	return p.finishedSteps

}

func readSingleStep(stepContent map[string]interface{}, p *parsingContext) {
	p.currentStep = new(stepDetails)
	p.nestingLevel = 0

	visitMap(stepContent, p)

	p.finishedSteps = append(p.finishedSteps, *p.currentStep)

}

func visitMap(stepContent map[string]interface{}, p *parsingContext) {

	for k, v := range stepContent {
		// fmt.Printf("Inside %s at %d\n", k, p.nestingLevel)

		switch v := v.(type) {
		case string:
			// fmt.Println(k, v, "(Found string at )", p.nestingLevel)
			storeStepInfo(k, v, p)
		case []interface{}:
			// fmt.Println(k, "(Found array):")
			storeSourcesInfo(k, v, p)
			// for i, u := range v {
			// 	fmt.Println("    ", i, u)
			// }
		case map[string]interface{}:
			// fmt.Println(k, "(Found object): ")
			p.nestingLevel++
			visitMap(v, p)
			p.nestingLevel--
			// default:
			// 	fmt.Println(k, v, "(Found unknown)")
		}

	}

}
func storeSourcesInfo(key string, values []interface{}, p *parsingContext) {
	if p.nestingLevel != 1 || key != "sources" || len(values) == 0 {
		return
	}
	p.currentStep.sourceURL = values[0].(string)
}

func storeStepInfo(key string, value string, p *parsingContext) {
	switch key {
	case "image":
		storeImageInfo(value, p)

	case "version":
		if p.nestingLevel == 1 {
			storeVersionInfo(value, p)
		}

	case "name":
		if p.nestingLevel != 1 {
			return
		}
		storeNameInfo(value, p)
	}
}

func storeImageInfo(fullDockerImage string, p *parsingContext) {
	var dockerImage dockerImageName

	if strings.Contains(fullDockerImage, ":") {
		imageAndTag := strings.Split(fullDockerImage, ":")
		dockerImage.baseImage = imageAndTag[0]
		dockerImage.hasTag = true
		dockerImage.tag = imageAndTag[1]

	} else {
		dockerImage.hasTag = false
		dockerImage.baseImage = fullDockerImage
	}
	p.currentStep.imagesUsed = append(p.currentStep.imagesUsed, dockerImage)
}

func storeVersionInfo(pluginVersion string, p *parsingContext) {
	p.currentStep.version = pluginVersion
}

func storeNameInfo(pluginName string, p *parsingContext) {
	p.currentStep.name = pluginName
}
