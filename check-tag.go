package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/genuinetools/reg/registry"
	"github.com/genuinetools/reg/repoutils"
	"github.com/sirupsen/logrus"
)

func connectToRegistryOfImage(imageAndTag *dockerImageName) *registry.Registry {
	image, err := registry.ParseImage(imageAndTag.BaseImage)
	if err != nil {
		log.Println(err)
		return nil
	}
	imageAndTag.Domain = image.Domain
	imageAndTag.Path = image.Path

	r, err := createRegistryClient(context.Background(), imageAndTag.Domain)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Connection successful")

	return r
}

func checkDockerImage(dockerHubConnection *registry.Registry, imageAndTag dockerImageName) bool {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Could not connect to Container registry, network error: %v \n", r)
		}
	}()

	tags, err := dockerHubConnection.Tags(context.Background(), imageAndTag.Path)
	if err != nil {
		log.Println(err)
		return false
	}
	sort.Strings(tags)

	log.Println("Found Dockerhub tags: ", tags)
	if !imageAndTag.HasTag {
		return true
	}

	for _, tag := range tags {
		if tag == imageAndTag.Tag {
			return true
		}
	}
	return false
}

func createRegistryClient(ctx context.Context, domain string) (*registry.Registry, error) {

	var (
		insecure    bool
		forceNonSSL bool

		timeout time.Duration

		authURL  string
		username string
		password string

		debug bool
	)

	// Use the auth-url domain if provided.
	authDomain := authURL
	if authDomain == "" {
		authDomain = domain
	}
	auth, err := repoutils.GetAuthConfig(username, password, authDomain)
	if err != nil {
		return nil, err
	}

	// Prevent non-ssl unless explicitly forced
	if !forceNonSSL && strings.HasPrefix(auth.ServerAddress, "http:") {
		return nil, fmt.Errorf("attempted to use insecure protocol! Use force-non-ssl option to force")
	}

	// Create the registry client.
	logrus.Infof("domain: %s", domain)
	logrus.Infof("server address: %s", auth.ServerAddress)
	return registry.New(ctx, auth, registry.Opt{
		Domain:   domain,
		Insecure: insecure,
		Debug:    debug,
		SkipPing: true, //Quay does not support ping https://github.com/genuinetools/reg/issues/198
		NonSSL:   forceNonSSL,
		Timeout:  timeout,
	})
}
