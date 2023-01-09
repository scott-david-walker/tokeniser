package main

import (
	"os"
	"testing"
)

func TestGetPackageVariable_WhenEmpty_ShouldReturnDefault(t *testing.T) {
	result, err := getPackageSpecificEnvironmentVariable("PREFIX", "#{")
	if result != "#{" || err != nil {
		t.Fatalf("Should have %q got %q, error: %v", "#{", result, err)
	}
}

func TestGetPackageVariable_WhenEmpty_AndNoDefault_ShouldReturnError(t *testing.T) {
	result, err := getPackageSpecificEnvironmentVariable("PREFIX", "")
	if err == nil {
		t.Fatalf("Should have received error but got %s", result)
	}
}

func TestGetPackageVariable_WhenNotEmpty_ShouldReturnValue(t *testing.T) {
	os.Setenv("INPUT_PREFIX", "value")
	result, err := getPackageSpecificEnvironmentVariable("PREFIX", "")
	if result != "value" || err != nil {
		t.Fatalf("Should have %q got %q, error: %v", "value", result, err)
	}
}

func TestGetEnvString_WhenEmpty_ShouldReturnError(t *testing.T) {
	result, err := getStringFromEnvironment("key")
	if err == nil {
		t.Fatalf("Should have recINPUT_PREFIXeived error but got %s", result)
	}
}

func TestGetEnvString_WhenNotEmpty_ShouldReturnValue(t *testing.T) {
	os.Setenv("key", "value")
	result, err := getStringFromEnvironment("key")
	if result != "value" || err != nil {
		t.Fatalf("Should have %q got %q, error: %v", "value", result, err)
	}
}
