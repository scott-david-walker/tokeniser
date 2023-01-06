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

func main() {
	fmt.Println(os.Environ())
	prefix, prefixErr := getPackageSpecificEnvironmentVariable("PREFIX", "#{")
	suffix, suffixErr := getPackageSpecificEnvironmentVariable("SUFFIX", "}#")
	failOnVariableNotFound := getPackageSpecificBoolEnvironmentVariable("FAIL_ON_VARIABLE_NOT_FOUND")
	glob, globError := getPackageSpecificEnvironmentVariable("FILES", "**")
	if prefixErr != nil {
		panic(prefixErr)
	}
	if suffixErr != nil {
		panic(suffixErr)
	}

	if globError != nil {
		panic(globError)
	}

	regex := buildRegexString(prefix, suffix)
	files := getFiles(glob)
	for _, file := range files {
		replaceValuesInFile(file, regex, failOnVariableNotFound)
	}
}

func getFiles(globPattern string) []string {
	glob := glob.MustCompile(globPattern)
	var files []string
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			println(err.Error())
			return nil
		}
		log.Println(path)
		match := glob.Match(path)
		if match && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}
func replaceValuesInFile(file string, regex *regexp.Regexp, failIfNotFound bool) {
	content, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		log.Fatal(readErr.Error())
		return
	}
	contentAsString := string(content[:])
	found := regex.FindAllString(contentAsString, -1)
	var elementMap = make(map[string]string)
	for _, f := range found {
		elementMap[f] = f
	}
	for k := range elementMap {
		val := k[2 : len(k)-2]
		envVal, envErr := getStringFromEnvironment(val)
		if envErr != nil && failIfNotFound {
			panic(errors.New(fmt.Sprintf("Replacable string %s found in file %s but has no corresponding replacement", val, file)))
		}
		contentAsString = strings.ReplaceAll(contentAsString, k, envVal)
	}
	writeErr := ioutil.WriteFile(file, []byte(contentAsString), 0)
	if writeErr != nil {
		panic(writeErr)
	}
}

func buildRegexString(prefix string, suffix string) *regexp.Regexp {
	regex := fmt.Sprintf("%s.*%s", prefix, suffix)
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
	variable := fmt.Sprintf("INPUT_%s", key)
	variable = os.Getenv(key)
	if variable == "" && defaultValue == "" {
		return variable, errors.New(fmt.Sprint("key with name %s is required", key))
	}

	if variable == "" {
		return defaultValue, nil
	}

	return variable, nil
}

// ignores

func getStringFromEnvironment(key string) (string, error) {
	variable := os.Getenv(key)
	if variable == "" {
		err := errors.New(fmt.Sprintf("key with name %s not found", key))
		return variable, err
	}

	return variable, nil
}
