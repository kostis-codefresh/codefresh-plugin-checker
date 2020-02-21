package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadSmallDataset(t *testing.T) {

	stepsFound := readJSON("testdata/input-small.json")

	assert.Len(t, stepsFound, 5, "Incorrect number of steps found")
}

func TestReadMediumDataset(t *testing.T) {

	stepsFound := readJSON("testdata/input-medium.json")

	assert.Equal(t, 14, len(stepsFound), "Incorrect number of steps found")
}

func TestReadFullDataset(t *testing.T) {

	stepsFound := readJSON("testdata/input-full.json")

	assert.Equal(t, 25, len(stepsFound), "Incorrect number of steps found")
}
