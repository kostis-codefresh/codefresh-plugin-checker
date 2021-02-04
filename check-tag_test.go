package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDockerImageWithValidImage(t *testing.T) {
	imageAndTag := dockerImageName{
		BaseImage: "codefresh/cf-sendmail",
		HasTag:    true,
		Tag:       "latest",
	}

	registryConnection := connectToRegistryOfImage(&imageAndTag)
	foundInRegistry := checkDockerImage(registryConnection, imageAndTag)

	assert.True(t, foundInRegistry, "Should be found in docker registry")
}

func TestCheckDockerImageWithFullDomain(t *testing.T) {
	imageAndTag := dockerImageName{
		BaseImage: "docker.io/codefresh/cf-sendmail",
		HasTag:    true,
		Tag:       "latest",
	}

	registryConnection := connectToRegistryOfImage(&imageAndTag)
	foundInRegistry := checkDockerImage(registryConnection, imageAndTag)

	assert.True(t, foundInRegistry, "Should be found in docker registry")
}

func TestCheckDockerImageWithInvalidImage(t *testing.T) {
	imageAndTag := dockerImageName{
		BaseImage: "foo",
		HasTag:    true,
		Tag:       "bar",
	}

	registryConnection := connectToRegistryOfImage(&imageAndTag)
	foundInRegistry := checkDockerImage(registryConnection, imageAndTag)

	assert.False(t, foundInRegistry, "Should be found in docker registry")
}

func TestCheckGCRImageWithValidImage(t *testing.T) {
	imageAndTag := dockerImageName{
		BaseImage: "gcr.io/cloud-builders/mvn",
		HasTag:    true,
		Tag:       "3.5.0-jdk-8",
	}

	registryConnection := connectToRegistryOfImage(&imageAndTag)
	foundInRegistry := checkDockerImage(registryConnection, imageAndTag)

	assert.True(t, foundInRegistry, "Should be found in docker registry")
}

func TestCheckGCRImageWithoutTag(t *testing.T) {
	imageAndTag := dockerImageName{
		BaseImage: "gcr.io/cloud-builders/git",
		HasTag:    false,
		Tag:       "",
	}

	registryConnection := connectToRegistryOfImage(&imageAndTag)
	foundInRegistry := checkDockerImage(registryConnection, imageAndTag)

	assert.True(t, foundInRegistry, "Should be found in docker registry")
}

func TestUnknownRegistry(t *testing.T) {
	imageAndTag := dockerImageName{
		BaseImage: "r.cfcr.io/jbadeau/gauge-typescript-plugin",
		HasTag:    false,
		Tag:       "",
	}

	registryConnection := connectToRegistryOfImage(&imageAndTag)
	foundInRegistry := checkDockerImage(registryConnection, imageAndTag)

	assert.False(t, foundInRegistry, "Should be found in docker registry")
}
