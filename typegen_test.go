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

		if schema.Properties == nil {
			log.Fatalln("schema properties cannot be empty")
		}
	}
}
