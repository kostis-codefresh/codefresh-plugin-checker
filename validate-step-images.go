package main

import (
	// "github.com/docker/distribution/digest"
	// "github.com/docker/distribution/manifest"
	// "github.com/docker/libtrust"
	//"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
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

	tmpl, err := template.ParseFiles("web/index.html.tpl")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("_site/index.html")
	if err != nil {
		log.Fatal("create file: ", err)
	}
	defer f.Close()

	tmpl.Execute(f, nil)

	if err != nil {
		log.Fatal("execute: ", err)
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
