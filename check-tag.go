package main

import (
	"fmt"
	"log"

	"github.com/heroku/docker-registry-client/registry"
)

func connectToDockerHub() *registry.Registry {
	url := "https://registry.hub.docker.com"
	username := "" // anonymous
	password := "" // anonymous
	dockeHubConnection, err := registry.New(url, username, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection successful")
	return dockeHubConnection
}

func checkDockerImage(dockerHubConnection *registry.Registry, imageAndTag dockerImageName) bool {
	tags, err := dockerHubConnection.Tags(imageAndTag.BaseImage)
	if err != nil {
		log.Println(err)
		return false
	}
	fmt.Println("Found dockerhub tags: ", tags)
	if imageAndTag.HasTag == false {
		return true
	}

	for _, tag := range tags {
		if tag == imageAndTag.Tag {
			return true
		}
	}
	return false

}
