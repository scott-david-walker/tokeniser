package main

import (
	"errors"
	"io/fs"
	"os"
	"regexp"
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
	defer genericPanic(t)
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

func TestGetConfiguration_ShouldPopulateObjectWithVariables(t *testing.T) {
	input := []string{"INPUT_PREFIX", "INPUT_SUFFIX", "INPUT_FILES"}

	for _, key := range input {
		os.Setenv(key, "test_"+key)
	}

	configuration := getConfiguration()

	if configuration.prefix != "test_INPUT_PREFIX" {
		t.Fatalf("Expected test_INPUT_PREFIX but received %s", configuration.prefix)
	}

	if configuration.suffix != "test_INPUT_SUFFIX" {
		t.Fatalf("Expected test_INPUT_SUFFIX but received %s", configuration.suffix)
	}

	if configuration.globPattern != "test_INPUT_FILES" {
		t.Fatalf("Expected test_INPUT_FILES but received %s", configuration.globPattern)
	}
}

func TestShouldPanicIfReadError(t *testing.T) {
	defer genericPanic(t)
	reg, err := regexp.Compile("#{.*?}#")
	if err != nil {
		panic("Regex should've compiled")
	}
	replaceValuesInFile("file", reg, configuration{}, mockRead("", true), nil)
}

func TestReplaceValues_WhenSomethingCanBeReplaced_ButHasNoEnvVariable_ShouldPanic(t *testing.T) {
	defer genericPanic(t)
	reg, err := regexp.Compile("#{.*?}#")
	if err != nil {
		panic("Regex should've compiled")
	}
	replaceValuesInFile("file", reg, configuration{}, mockRead("#{blah}#", false), nil)
}

func TestReplaceValues_WhenSomethingCanBeReplaced_ShouldReplace(t *testing.T) {
	os.Setenv("blah", "value")
	reg, err := regexp.Compile("#{.*?}#")
	if err != nil {
		panic("Regex should've compiled")
	}

	writeFile := func(filename string, data []byte, perm fs.FileMode) error {
		content := string(data[:])
		if content != "value" {
			t.Fatalf("Expected %s but received %s", "value", content)
		}
		return nil
	}
	replaceValuesInFile("file", reg, configuration{prefix: "#{", suffix: "}#"}, mockRead("#{blah}#", false), writeFile)
}

func mockRead(content string, shouldPanic bool) func(filename string) ([]byte, error) {
	return func(filename string) ([]byte, error) {
		var err error = nil
		bytes := []byte(content)
		if shouldPanic {
			err = errors.New("error")
		}
		return bytes, err
	}
}

func genericPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}
