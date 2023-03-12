package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	TSString  = "string"
	TSNumber  = "number"
	TSBoolean = "boolean"
	TSArray   = "Array<%s>"
	TSObject  = "object"
	TSAny     = "any"
)

type JSONSchema struct {
	ID          string                 `json:"$id,omitempty"`
	Schema      string                 `json:"$schema,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Definitions map[string]*JSONSchema `json:"definitions,omitempty"`
}

type TSField struct {
	Name     string
	Type     string
	Required bool
}

type TSInterface struct {
	Name   string
	Fields []TSField
}

func main() {
	usage := `usage: codegen <path/to/schema.json>`

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	input := os.Args[1]

	buff, err := os.ReadFile(input)

	if err != nil {
		log.Fatal(err)
	}

	schema, err := Parse(buff)

	if err != nil {
		log.Fatal(err)
	}

	tsInterface, err := schema.GenerateInterface()

	if err != nil {
		log.Fatal(err)
	}

	file := strings.Split(input, "/")[len(strings.Split(input, "/"))-1]

	dts := strings.Split(file, ".schema.json")[0] + ".d.ts"

	err = os.WriteFile(dts, []byte(tsInterface), 0644)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("generated types for %s schema in %s\n", schema.Title, dts)
}

func Parse(raw []byte) (JSONSchema, error) {
	var schema JSONSchema

	if raw == nil {
		return schema, fmt.Errorf("raw data cannot be empty")
	}

	err := json.Unmarshal(raw, &schema)

	if err != nil {
		return schema, err
	}

	return schema, nil
}

func (schema JSONSchema) GenerateInterface() (string, error) {
	var tsInterface TSInterface

	if schema.Title != "" {
		tsInterface.Name = schema.Title
	}

	for k, v := range schema.Properties {
		var tsField TSField

		tt := v.(map[string]interface{})["type"]

		tsField.Name = k

		switch tt {
		case "string":
			tsField.Type = TSString
		case "integer":
			tsField.Type = TSNumber
		case "number":
			tsField.Type = TSNumber
		case "boolean":
			tsField.Type = TSBoolean
		case "array":
			// TODO: handle array types
			tsField.Type = fmt.Sprintf(TSArray, TSString)
		case "object":
			tsField.Type = TSObject
		default:
			return "", fmt.Errorf("unknown type %s", tt)
		}
		tsInterface.Fields = append(tsInterface.Fields, tsField)
	}

	for _, v := range schema.Required {
		// fmt.Printf("required field: %s\n", v)
		for i, j := range tsInterface.Fields {
			if j.Name == v {
				tsInterface.Fields[i].Required = true
			}
		}
	}
	return tsInterface.String(), nil
}

func (inter TSInterface) String() string {
	var out string

	out += fmt.Sprintf("export interface %s {\n", inter.PascalCase())

	for _, v := range inter.Fields {
		name := strings.ToLower(v.PascalCase()[:1]) + v.PascalCase()[1:]

		if !v.Required {
			out += fmt.Sprintf("\t%s?: %s;\n", name, v.Type)
		} else {
			out += fmt.Sprintf("\t%s: %s;\n", name, v.Type)
		}
	}

	out += "}\n"

	return out
}

func (field TSField) PascalCase() string {
	cc := strings.Split(field.Name, "_")

	var out string

	caser := cases.Title(language.English)

	for _, v := range cc {
		out += caser.String(v)
	}

	return out
}

func (inter TSInterface) PascalCase() string {
	cc := strings.Split(inter.Name, "_")

	var out string

	caser := cases.Title(language.English)

	for _, v := range cc {
		out += caser.String(v)
	}

	return out
}
