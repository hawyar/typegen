package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var dir = "examples"
	var files []string

	schemas, err := os.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	for _, f := range schemas {
		if f.IsDir() {
			continue
		}
		files = append(files, f.Name())
	}

	for _, f := range files {
		buff, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, f))

		if err != nil {
			panic(err)
		}

		schema, err := Parse(buff)

		if err != nil {
			panic(err)
		}

		if schema.Title == "" {
			log.Fatalln("schema title cannot be empty")
		}
	}
}

func TestParse(t *testing.T) {
	example := "./example/account.schema.json"

	buff, err := os.ReadFile(example)

	if err != nil {
		t.Fatal(err)
	}

	schema, err := Parse(buff)

	if err != nil {
		t.Fatal(err)
	}

	if schema.Title == "" {
		t.Fatal("schema title cannot be empty")
	}
}
