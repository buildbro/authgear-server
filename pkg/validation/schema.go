package validation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
	"github.com/iawaknahc/jsonschema/pkg/jsonschema"
	jsonschemaformat "github.com/iawaknahc/jsonschema/pkg/jsonschema/format"
)

type SimpleSchema struct {
	col *jsonschema.Collection
}

func NewSimpleSchema(schema string) *SimpleSchema {
	col := jsonschema.NewCollection()
	col.AddSchema(strings.NewReader(schema), "")
	return &SimpleSchema{
		col: col,
	}
}

func (s *SimpleSchema) RegisterFormat(format string, checker jsonschemaformat.FormatChecker) {
	s.col.FormatChecker[format] = checker
}

func (s *SimpleSchema) Validator() *SchemaValidator {
	return &SchemaValidator{Schema: s.col}
}

type MultipartSchema struct {
	mainPartID string
	parts      map[string]interface{}
	col        *jsonschema.Collection
}

func NewMultipartSchema(mainPartID string) *MultipartSchema {
	return &MultipartSchema{
		mainPartID: mainPartID,
		parts:      map[string]interface{}{},
		col:        nil,
	}
}

func (s *MultipartSchema) Add(partID string, schema string) *MultipartSchema {
	if s.col != nil {
		panic("validation: cannot add part when schema is already instantiated")
	}
	var schemaObj interface{}
	if err := json.Unmarshal([]byte(schema), &schemaObj); err != nil {
		panic(fmt.Sprintf("validation: invalid schema part '%s': %s", partID, err))
	}
	s.parts[partID] = schemaObj
	return s
}

func (s *MultipartSchema) DumpSchemaString(pretty bool) (schemaString string, err error) {
	schema := map[string]interface{}{
		"$defs": s.parts,
		"$ref":  jsonpointer.T([]string{"$defs", s.mainPartID}),
	}

	var schemaJSON []byte
	if pretty {
		schemaJSON, err = json.MarshalIndent(schema, "", "  ")
	} else {
		schemaJSON, err = json.Marshal(schema)
	}
	if err != nil {
		return
	}

	schemaString = string(schemaJSON)
	return
}

func (s *MultipartSchema) Instantiate() *MultipartSchema {
	if _, ok := s.parts[s.mainPartID]; !ok {
		panic(fmt.Sprintf("validaiton: main part '%s' is not added", s.mainPartID))
	}

	schemaString, err := s.DumpSchemaString(false)
	if err != nil {
		panic("validation: invalid JSON schema: " + err.Error())
	}

	// Do not forget the parts so that we can dump the schema later.
	// s.parts = nil
	s.col = jsonschema.NewCollection()
	s.col.AddSchema(strings.NewReader(schemaString), "")
	return s
}

func (s *MultipartSchema) RegisterFormat(format string, checker jsonschemaformat.FormatChecker) {
	if s.col == nil {
		panic("validation: JSON schema is not instantiated")
	}
	s.col.FormatChecker[format] = checker
}

func (s *MultipartSchema) Validator() *SchemaValidator {
	if s.col == nil {
		panic("validation: JSON schema is not instantiated")
	}
	return &SchemaValidator{Schema: s.col}
}

func (s *MultipartSchema) PartValidator(partID string) *SchemaValidator {
	if s.col == nil {
		panic("validation: JSON schema is not instantiated")
	}
	return &SchemaValidator{
		Schema:    s.col,
		Reference: jsonpointer.T([]string{"$defs", partID}).Fragment(),
	}
}
