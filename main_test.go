package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	expected := "Hello, Go!"
	result := getHello()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

// Refactor the main function to make it testable
func getHello() string {
	return "Hello, Go!"
}
