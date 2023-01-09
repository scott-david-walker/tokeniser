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

func TestGetPackageBoolVariable_WhenNotEmpty_ShouldReturnValue(t *testing.T) {
	input := []string{"true", "false"}
	expectation := []bool{true, false}
	for index, val := range input {
		os.Setenv("INPUT_KEY", val)
		result := getPackageSpecificBoolEnvironmentVariable("KEY")
		if result != expectation[index] {
			t.Fatalf("Should have %t got %t", expectation[index], result)
		}
	}
}

func TestGetPackageBoolVariable_WhenEmpty_ShouldReturnTrue(t *testing.T) {
	os.Setenv("INPUT_KEY", "")
	result := getPackageSpecificBoolEnvironmentVariable("KEY")
	if result != true {
		t.Fatalf("Should have true got %t", result)
	}
}

func TestGetPackageBoolVariable_WhenValueIsNotParseableAsBool_ShouldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	os.Setenv("INPUT_KEY", "panic")
	_ = getPackageSpecificBoolEnvironmentVariable("KEY")
}

func TestGetEnvString_WhenNotEmpty_ShouldReturnValue(t *testing.T) {
	os.Setenv("key", "value")
	result, err := getStringFromEnvironment("key")
	if result != "value" || err != nil {
		t.Fatalf("Should have %q got %q, error: %v", "value", result, err)
	}
}

func TestRegexForFindingSearchString_ShouldMatch(t *testing.T) {
	input := "#{search}#"
	reg := buildRegexString("#{", "}#")
	found := reg.FindAllString(input, -1)
	numOfFound := 0
	for range found {
		numOfFound++
	}

	if numOfFound != 1 {
		t.Fatalf("Expected 1 match but found %d", numOfFound)
	}
}

func TestRegexForFindingSearchString_WhenMultipleMatchesOnSameLine_ShouldMatch(t *testing.T) {
	input := "#{search}##{another}#"
	reg := buildRegexString("#{", "}#")
	found := reg.FindAllString(input, -1)
	numOfFound := 0
	for range found {
		numOfFound++
	}

	if numOfFound != 2 {
		t.Fatalf("Expected 2 matches but found %d", numOfFound)
	}
}
