package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("Checking node image")

	if len(os.Args) < 2 {
		fmt.Println("Input JSON file is missing.")
		os.Exit(1)
	}

	err := os.MkdirAll("_site", 0755)
	if err != nil {
		log.Fatal(err)
	}

	copy("web/style.css", "_site/style.css")

	stepsFound := readJSON(os.Args[1])

	fmt.Printf("Found %d steps\n", len(stepsFound))

	postProcessSteps(stepsFound)

	tmpl, err := template.ParseFiles("web/index.html.tpl")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("_site/index.html")
	if err != nil {
		log.Fatal("create file: ", err)
	}
	defer f.Close()

	type templateData struct {
		Now           time.Time
		FinishedSteps []stepDetails
	}

	tData := templateData{
		Now:           time.Now(),
		FinishedSteps: stepsFound,
	}

	err = tmpl.Execute(f, tData)

	if err != nil {
		log.Fatal("execute: ", err)
	}

}

func postProcessSteps(stepsFound []stepDetails) {

	dockeHubConnection := connectToDockerHub()

	for i := range stepsFound {
		if len(stepsFound[i].ImagesUsed) == 0 {
			stepsFound[i].Status = incomplete
			continue
		}

		stepsFound[i].Status = ok

		for y := range stepsFound[i].ImagesUsed {
			stepsFound[i].ImagesUsed[y].FoundInRegistry = checkDockerImage(dockeHubConnection, stepsFound[i].ImagesUsed[y])
			if stepsFound[i].ImagesUsed[y].FoundInRegistry == false {
				stepsFound[i].Status = notOk
			}
		}

	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copy(src string, dst string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	checkErr(err)
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	checkErr(err)
}
