package main

import (
	"testing"
)

func TestGetPrefix_WhenEmpty_ShouldReturnDefault(t *testing.T) {
	result, err := GetPackageSpecificEnvironmentVariable("PREFIX", "#{")
	if result != "#{" || err != nil {
		t.Fatalf("Should have %q got %q, error: %v", "#{", result, err)
	}
}
