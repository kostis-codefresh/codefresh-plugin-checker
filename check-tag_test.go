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

	dockeHubConnection := connectToDockerHub()
	foundInRegistry := checkDockerImage(dockeHubConnection, imageAndTag)

	assert.True(t, foundInRegistry, "Should be found in docker registry")
}

func TestCheckDockerImageWithInvalidImage(t *testing.T) {
	imageAndTag := dockerImageName{
		BaseImage: "foo",
		HasTag:    true,
		Tag:       "bar",
	}

	dockeHubConnection := connectToDockerHub()
	foundInRegistry := checkDockerImage(dockeHubConnection, imageAndTag)

	assert.False(t, foundInRegistry, "Should be found in docker registry")
}

func TestPrepareBaseImageWithoutPrefixes(t *testing.T) {
	output := prepareBaseImage("codefresh/cf-sendmail")

	assert.EqualValues(t, "codefresh/cf-sendmail", output, "Should be the same")
}

func TestPrepareBaseImageWithDockerioPrefix(t *testing.T) {
	output := prepareBaseImage("docker.io/codefresh/cf-sendmail")

	assert.EqualValues(t, "codefresh/cf-sendmail", output, "Should remove 'docker.io/'")
}
