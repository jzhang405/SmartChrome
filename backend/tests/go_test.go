package tests

import (
	"testing"
)

func TestInit(t *testing.T) {
	// Basic test to ensure testing framework is working
	if true != true {
		t.Error("Test framework not working")
	}
}