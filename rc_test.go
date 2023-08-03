package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

func checkExpectedResults(path string, testConfig configFile, expected []illegal, t *testing.T) {
	config = testConfig
	restrictLiterals()
	allowBuiltins()
	allowCasting()
	allowFunctions()
	// Reset globals after running
	defer func() {
		config = configFile{}
		noLitRegexp = nil
		allowedFun = make(map[string]map[string]bool)
		allowedRep = make(map[string]int)
	}()

	load := make(loadedSource)
	err := loadProgram(filepath.Dir(path), load)
	if err != nil {
		t.Log("Failed to parse source:", err)
		t.Fail()
	}

	currentPath := filepath.Dir(path)

	info := analyzeProgram(path, currentPath, load)
	if len(expected) == 0 && len(info.illegals) == 0 {
		// Dirty hack for empty slices
	} else if !reflect.DeepEqual(expected, info.illegals) {
		t.Log("Expected:", expected)
		t.Log("Got:", info.illegals)
		t.FailNow()
	}
	fmt.Printf("%s: Ok!\n", t.Name())
}

func TestAllowedFunctionsFail(t *testing.T) {
	testConfig := configFile{}

	expected := []illegal{
		illegal{
			T:    "illegal-import",
			Name: "fmt",
			Pos:  "tests/allowedFunctions/main.go:3:8",
		},
		illegal{
			T:    "illegal-access",
			Name: "fmt.Println",
			Pos:  "tests/allowedFunctions/main.go:6:2",
		},
		illegal{
			T:    "illegal-definition",
			Name: "main",
			Pos:  "tests/allowedFunctions/main.go:5:1",
		},
	}
	checkExpectedResults("tests/allowedFunctions/main.go", testConfig, expected, t)
}

func TestAllowedFunctions(t *testing.T) {
	testConfig := configFile{
		AllowedFunctions: []string{"fmt.*"},
	}

	expected := []illegal{}
	checkExpectedResults("tests/allowedFunctions/main.go", testConfig, expected, t)
}
