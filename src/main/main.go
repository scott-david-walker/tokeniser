package main

import (
	"errors"
	"fmt"
	"github.com/gobwas/glob"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type configuration struct {
	prefix                 string
	suffix                 string
	globPattern            string
	failOnVariableNotFound bool
}

func main() {
	configuration := getConfiguration()
	regex := buildRegexString(configuration.prefix, configuration.suffix)
	files := getFiles(configuration.globPattern)
	for _, file := range files {
		replaceValuesInFile(file, regex, configuration)
	}
}

func getConfiguration() configuration {
	const prefixKey = "PREFIX"
	const suffixKey = "SUFFIX"
	const filesKey = "FILES"
	const failOnNotFoundKey = "FAIL-IF-NO-PROVIDED-REPLACEMENT"
	prefix, prefixErr := getPackageSpecificEnvironmentVariable(prefixKey, "#{")
	suffix, suffixErr := getPackageSpecificEnvironmentVariable(suffixKey, "}#")
	globPattern, globError := getPackageSpecificEnvironmentVariable(filesKey, "**")
	failOnVariableNotFound := getPackageSpecificBoolEnvironmentVariable(failOnNotFoundKey)

	if prefixErr != nil {
		panic(prefixErr)
	}
	if suffixErr != nil {
		panic(suffixErr)
	}

	if globError != nil {
		panic(globError)
	}

	return configuration{
		prefix:                 prefix,
		suffix:                 suffix,
		globPattern:            globPattern,
		failOnVariableNotFound: failOnVariableNotFound,
	}
}

func getFiles(globPattern string) []string {
	compiledGlob := glob.MustCompile(globPattern)
	var files []string
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err.Error())
			return nil
		}
		match := compiledGlob.Match(path)
		if match && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func replaceValuesInFile(file string, regex *regexp.Regexp, config configuration) {
	content, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		panic(readErr)
	}
	contentAsString := string(content[:])
	found := regex.FindAllString(contentAsString, -1)
	if len(found) == 0 {
		return
	}
	var elementMap = make(map[string]string)
	for _, foundItem := range found {
		elementMap[foundItem] = foundItem
	}
	for key := range elementMap {
		val := key[len(config.prefix) : len(key)-len(config.suffix)]
		envVal, envErr := getStringFromEnvironment(val)
		if envErr != nil && config.failOnVariableNotFound {
			panic(errors.New(fmt.Sprintf("Replacable string '%s' found in file '%s' but has no corresponding replacement", val, file)))
		}
		log.Println(fmt.Sprintf("Replacing value '%s' with '%s' in file '%s'", key, envVal, file))
		contentAsString = strings.ReplaceAll(contentAsString, key, envVal)
	}
	writeErr := ioutil.WriteFile(file, []byte(contentAsString), 0)
	if writeErr != nil {
		panic(writeErr)
	}
}

func buildRegexString(prefix string, suffix string) *regexp.Regexp {
	regex := fmt.Sprintf("%s.*?%s", prefix, suffix)
	reg, err := regexp.Compile(regex)
	if err != nil {
		panic(errors.New(err.Error()))
	}
	return reg
}

func getPackageSpecificBoolEnvironmentVariable(key string) bool {
	v, err := getPackageSpecificEnvironmentVariable(key, "")

	if err != nil {
		return true
	}
	bv, bErr := strconv.ParseBool(v)

	if bErr != nil {
		panic(errors.New(fmt.Sprintf("cannot convert %s to boolean", v)))
	}
	return bv
}

func getPackageSpecificEnvironmentVariable(key string, defaultValue string) (string, error) {
	newKey := fmt.Sprintf("INPUT_%s", key)
	variable := os.Getenv(newKey)
	if variable == "" && defaultValue == "" {
		return variable, errors.New(fmt.Sprintf("key with name %s is required", key))
	}

	if variable == "" {
		return defaultValue, nil
	}

	return variable, nil
}

func getStringFromEnvironment(key string) (string, error) {
	variable := os.Getenv(key)
	if variable == "" {
		err := errors.New(fmt.Sprintf("key with name %s not found", key))
		return variable, err
	}

	return variable, nil
}
