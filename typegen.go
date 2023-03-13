package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
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

	JSONString  = "string"
	JSONInteger = "integer"
	JSONNumber  = "number"
	JSONBool    = "boolean"
	JSONArray   = "array"
	JSONObject  = "object"
	JSNULL      = "null"
)

type JSONSchemaItem struct {
	Type        string                     `json:"type"`
	Items       *JSONSchemaItem            `json:"items"`
	Properties  map[string]*JSONSchemaItem `json:"properties"`
	Required    []string                   `json:"required"`
	Description string                     `json:"description"`
}
type JSONSchema struct {
	Identifier string                     `json:"$id,omitempty"`
	Schema     string                     `json:"$schema,omitempty"`
	Title      string                     `json:"title,omitempty"`
	Type       string                     `json:"type,omitempty"`
	Properties map[string]*JSONSchemaItem `json:"properties,omitempty"`
	Required   []string                   `json:"required,omitempty"`
}

type TSField struct {
	Name      string
	Type      string
	Required  bool
	Reference map[string]*TSInterface
}

type TSInterface struct {
	Name       string
	Fields     []TSField
	References map[string]*TSInterface
}

func main() {
	var inputFile, outputDir string

	usage := `usage: codegen <path/to/schema.json>`

	if len(os.Args) < 2 {
		log.Fatal(usage)
	}

	if len(os.Args) == 3 {
		inputFile = os.Args[1]
		outputDir = os.Args[2]
	} else {
		inputFile = os.Args[1]
	}

	if inputFile == "" {
		log.Fatalln("empty input file path: specify a valid path to a JSON Schema file")
	}

	if outputDir == "" {
		curr, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
		}

		outputDir = curr
	}

	raw, err := os.ReadFile(inputFile)

	if err != nil {
		log.Fatal(err)
	}

	schema, err := Parse(raw)

	if err != nil {
		log.Fatalf("failed to parse JSON Schema %s: %s\n", inputFile, err)
	}

	root, err := schema.GenerateInterface()

	if err != nil {
		log.Fatalf("failed to generate TS types for %s: %s\n", schema.Title, err)
	}

	outputFileName := path.Join(outputDir, PascalCase(schema.Title)+".d.ts")

	err = Write(outputFileName, root)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated TS types for %s saved in %s\n", schema.Title, outputFileName)
}

func Write(file string, roots TSInterface) error {
	var out string

	note := fmt.Sprintf("// Generated TS types for %s JSON schema\n\n", roots.Name)

	out += note

	out += roots.String()

	err := os.WriteFile(file, []byte(out), 0644)

	if err != nil {
		return err
	}
	return nil
}

func (item *JSONSchemaItem) UnmarshalJSON(data []byte) error {
	var temp struct {
		Type       string                     `json:"type"`
		Items      *JSONSchemaItem            `json:"items"`
		Properties map[string]*JSONSchemaItem `json:"properties"`
		Required   []string                   `json:"required"`
		Desc       string                     `json:"description"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	item.Type = temp.Type
	item.Items = temp.Items
	item.Properties = temp.Properties
	item.Required = temp.Required
	item.Description = temp.Desc

	return nil
}

func Parse(raw []byte) (*JSONSchema, error) {
	if raw == nil {
		return nil, fmt.Errorf("raw data cannot be empty")
	}

	var schema JSONSchema

	err := json.Unmarshal(raw, &schema)

	if err != nil {
		return nil, err
	}

	return &schema, nil
}

func (schema JSONSchema) GenerateInterface() (TSInterface, error) {
	var root TSInterface

	refs := make(map[string]*TSInterface)

	if schema.Title == "" {
		return root, fmt.Errorf("schema title cannot be empty")
	}

	if schema.Properties == nil {
		return root, fmt.Errorf("schema has no properties")
	}

	root.Name = schema.Title

	for k, v := range schema.Properties {
		var tsField TSField

		tsField.Name = k

		switch v.Type {
		case "string":
			tsField.Type = TSString
		case "integer":
			tsField.Type = TSNumber
		case "number":
			tsField.Type = TSNumber
		case "boolean":
			tsField.Type = TSBoolean
		case "array":
			if v.Items == nil {
				tsField.Type = TSAny
				root.Fields = append(root.Fields, tsField)
				break
			}

			if v.Items.Type == "object" {
				if v.Items.Properties == nil {
					tsField.Type = TSAny
					root.Fields = append(root.Fields, tsField)
					break
				}

				var ref TSInterface

				ref.Name = k

				if ref.Name[len(ref.Name)-1] == 's' {
					ref.Name = strings.TrimSuffix(ref.Name, "s")
				}

				ref.Name = PascalCase(ref.Name)

				for i, j := range v.Items.Properties {
					var refField TSField

					refField.Name = i

					switch j.Type {
					case "string":
						refField.Type = TSString
					case "integer":
						refField.Type = TSNumber
					case "number":
						refField.Type = TSNumber
					case "boolean":
						refField.Type = TSBoolean
					default:
						return root, fmt.Errorf("unknown type %s", j.Type)
					}

					for _, l := range v.Items.Required {
						if l == i {
							refField.Required = true
						}
					}

					ref.Fields = append(ref.Fields, refField)
					tsField.Type = fmt.Sprintf(TSArray, ref.Name)
					refs[ref.Name] = &ref
				}
			}
		case "object":
			tsField.Type = TSObject
		default:
			return root, fmt.Errorf("unknown type %s", v.Type)
		}

		root.Fields = append(root.Fields, tsField)
		root.References = refs
	}

	for _, v := range schema.Required {
		for i, j := range root.Fields {
			if j.Name == v {
				root.Fields[i].Required = true
			}
		}
	}

	return root, nil
}

func (root TSInterface) String() string {
	var out string

	for _, v := range root.References {
		out += v.String()
	}

	out += fmt.Sprintf("export interface %s {\n", PascalCase(root.Name))

	for _, v := range root.Fields {
		name := strings.ToLower(PascalCase(v.Name)[:1]) + PascalCase(v.Name)[1:]

		if !v.Required {
			out += fmt.Sprintf("\t%s?: %s;\n", name, v.Type)
		} else {
			out += fmt.Sprintf("\t%s: %s;\n", name, v.Type)
		}
	}

	out += "}\n"

	return out
}

func PascalCase(str string) string {
	cc := strings.Split(str, "_")

	var out string

	caser := cases.Title(language.English)

	for _, v := range cc {
		out += caser.String(v)
	}

	return out
}
