package main

import (
	"testing"
)

// TestCreateFlair tests if flair is there is any error
// while creating the flair
func TestCreateFlair(t *testing.T) {

	testingBit = true
	err := PrepareTemplate()
	if err != nil {
		t.Error(err)
	}
	_, err = CreateFlair("arpitjindal97", "clean")

	if err != nil {
		t.Error(err)
	}

}
