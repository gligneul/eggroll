// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package compiler

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

type tmplData struct {
	Package  string
	JsonAbi  string
	Structs  []*tmplMessageSchema
	Schemas  []*tmplMessageSchema
	Advances []*tmplMessageSchema
	Inspects []*tmplMessageSchema
}

type tmplMessageSchema struct {
	Kind   string
	Doc    string
	GoName string
	ID     string
	Fields []tmplFieldSchema
}

type tmplFieldSchema struct {
	Kind   string
	Doc    string
	GoName string
	Type   string
}

// Generate the EggRoll Go binding for the ast.
func generateGo(ast astSchema, packageName string) []byte {
	var data tmplData
	data.Package = packageName
	data.JsonAbi = string(generateAbi(ast))
	for _, struct_ := range ast.Structs {
		schema := generateTmplMessage(struct_, ast.Structs)
		data.Structs = append(data.Structs, &schema)
	}
	for _, report := range ast.Reports {
		schema := generateTmplMessage(report, ast.Structs)
		data.Structs = append(data.Structs, &schema)
		data.Schemas = append(data.Schemas, &schema)
	}
	for _, advance := range ast.Advances {
		schema := generateTmplMessage(advance, ast.Structs)
		data.Structs = append(data.Structs, &schema)
		data.Schemas = append(data.Schemas, &schema)
		data.Advances = append(data.Advances, &schema)
	}
	for _, inspect := range ast.Inspects {
		schema := generateTmplMessage(inspect, ast.Structs)
		data.Structs = append(data.Structs, &schema)
		data.Schemas = append(data.Schemas, &schema)
		data.Inspects = append(data.Inspects, &schema)
	}

	// generate code using template
	tmpl := template.Must(template.New("eggroll").Parse(tmplSource))
	var codeBuffer bytes.Buffer
	err := tmpl.Execute(&codeBuffer, data)
	if err != nil {
		panic(err)
	}

	// format the source code
	code, err := format.Source(codeBuffer.Bytes())
	if err != nil {
		panic(fmt.Errorf("%v\n%v", err, codeBuffer.String()))
	}
	return code
}

// Generate a template schema from the message.
func generateTmplMessage(message messageSchema, structs []messageSchema) tmplMessageSchema {
	var tmplMessage tmplMessageSchema
	tmplMessage.Kind = message.Name
	tmplMessage.Doc = generateDoc(message.Doc)
	tmplMessage.GoName = captalize(message.Name)
	tmplMessage.ID = captalize(message.Name) + "ID"
	for _, field := range message.Fields {
		var tmplField tmplFieldSchema
		tmplField.Kind = field.Name
		tmplField.Doc = generateDoc(field.Doc)
		tmplField.GoName = captalize(field.Name)
		tmplField.Type = generateGoType(field.type_, structs)
		tmplMessage.Fields = append(tmplMessage.Fields, tmplField)
	}
	return tmplMessage
}

// Generate a Go type.
func generateGoType(type_ any, structs []messageSchema) string {
	switch type_ := type_.(type) {
	case typeBool:
		return "bool"
	case typeInt:
		prefix := ""
		if !type_.Signed {
			prefix = "u"
		}
		switch type_.Bits {
		case 8, 16, 32, 64:
			return fmt.Sprintf("%vint%v", prefix, type_.Bits)
		}
		return "*big.Int"
	case typeAddress:
		return "common.Address"
	case typeBytes:
		return "[]byte"
	case typeString:
		return "string"
	case typeArray:
		return "[]" + generateGoType(type_.Elem, structs)
	case typeStructRef:
		struct_ := structs[type_.Index]
		return captalize(struct_.Name)
	default:
		// This should not happen
		panic(fmt.Errorf("invalid type: %T", type_))
	}
}

// Prefix each line of the doc string with //
func generateDoc(doc string) string {
	if doc == "" {
		return doc
	}
	if doc[len(doc)-1] == '\n' {
		// remove trailing \n
		doc = doc[:len(doc)-1]
	}
	lines := strings.Split(doc, "\n")
	for i, line := range lines {
		lines[i] = "// " + line
	}
	return strings.Join(lines, "\n")
}

// Captalize the first letter.
func captalize(name string) string {
	return strings.ToUpper(name[0:1]) + name[1:]
}

const tmplSource = `// Code generated by EggRoll - DO NOT EDIT.

package {{.Package}}

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll/pkg/eggtypes"
	"github.com/gligneul/eggroll/pkg/eggroll"
)

var (
	_ = big.NewInt
	_ = common.Big1
	_ = eggtypes.MustAddSchema
)


// Messages encoded as JSON ABI.
const _JSON_ABI = ` + "`" + `{{.JsonAbi}}
` + "`" + `

// Solidity ABI.
var _abi abi.ABI

//
// Struct Types
//

{{range $struct := .Structs}}
	{{- $struct.Doc}}
	type {{$struct.GoName}} struct {
	{{range $field := .Fields}}
		{{- $field.Doc}}
		{{$field.GoName}} {{$field.Type}}{{end}}
	}
{{end}}

//
// ID for each schema
//

{{range $schema := .Schemas}}
	// 4-byte function selector of {{$schema.Kind}}
	var {{$schema.ID}} eggtypes.ID
{{end}}

//
// Encode functions for each message schema
//

{{range $schema := .Schemas}}
	// Encode {{$schema.Kind}} into binary data.
	func Encode{{$schema.GoName}}(
		{{- range $field := .Fields}}
			{{$field.GoName}} {{$field.Type}},
		{{- end}}
	) []byte {
		values := make([]any, {{- len $schema.Fields}})
		{{- range $i, $field := .Fields}}
			values[{{$i}}] = {{$field.GoName}}
		{{- end}}
		data, err := _abi.Methods["{{$schema.Kind}}"].Inputs.PackValues(values)
		if err != nil {
			panic(fmt.Sprintf("failed to encode {{$schema.Kind}}: %v", err))
		}
		return append({{$schema.ID}}[:], data...)
	}

	// Encode {{$schema.Kind}} into binary data.
	func (v {{$schema.GoName}}) Encode() []byte {
		return Encode{{$schema.GoName}}(
		{{- range $field := .Fields}}
			v.{{$field.GoName}},
		{{- end}}
		)
	}
{{end}}

//
// Decode functions for each message schema
//

{{range $schema := .Schemas}}
	func _decode_{{$schema.GoName}}(values []any) (any, error) {
		if len(values) != {{len $schema.Fields}} {
			return nil, fmt.Errorf("wrong number of values")
		}
		{{- if $schema.Fields}}
			var ok bool
		{{- end}}
		var v {{$schema.GoName}}
		{{- range $i, $field := .Fields}}
			v.{{$field.GoName}}, ok = values[{{$i}}].({{$field.Type}})
			if !ok {
				return nil, fmt.Errorf("failed to decode {{$schema.Kind}}.{{$field.Kind}}")
			}
		{{- end}}
		return v, nil
	}
{{end}}

//
// Init function
//

func init() {
	var err error
	_abi, err = abi.JSON(strings.NewReader(_JSON_ABI))
	if err != nil {
		// This should not happen
		panic(fmt.Sprintf("failed to decode ABI: %v", err))
	}
	{{- range $schema := .Schemas}}
		{{$schema.ID}} = eggtypes.ID(_abi.Methods["{{$schema.Kind}}"].ID)
		eggtypes.MustAddSchema(eggtypes.MessageSchema{
			ID:        {{$schema.ID}},
			Kind:      "{{$schema.Kind}}",
			Arguments: _abi.Methods["{{$schema.Kind}}"].Inputs,
			Decoder:   _decode_{{$schema.GoName}},
		})
	{{- end}}
}

//
// Middleware
//

// High-level contract
type iContract interface {
	{{range $advance := .Advances}}
		{{$advance.Doc}}
		{{$advance.GoName}}(
			eggroll.Env,
			{{- range $field := $advance.Fields}}
				{{$field.Type}},
			{{- end}}
		) error
	{{end}}

	{{range $inspect := .Inspects}}
		{{$inspect.Doc}}
		{{$inspect.GoName}}(
			eggroll.EnvReader,
			{{- range $field := $inspect.Fields}}
				{{$field.Type}},
			{{- end}}
		) error
	{{end}}
}

// Middleware that implements the EggRoll Middleware interface.
// The middleware requires a high-level contract to work.
type Middleware struct {
	contract iContract
}

func (m Middleware) Advance(env eggroll.Env, input []byte) error {
	{{- if .Advances}}
		unpacked, err := eggtypes.Decode(input)
		if err != nil {
			return err
		}
		env.Logf("middleware: received %#v", unpacked)
		switch input := unpacked.(type) {
		{{- range $advance := .Advances}}
		case {{$advance.GoName}}:
			return m.contract.{{$advance.GoName}}(
				env,
				{{- range $field := $advance.Fields}}
					input.{{$field.GoName}},
				{{- end}}
			)
		{{- end}}
		default:
			return fmt.Errorf("middleware: input isn't an advance")
		}
	{{- else}}
		return fmt.Errorf("advance not supported")
	{{- end}}
}

func (m Middleware) Inspect(env eggroll.EnvReader, input []byte) error {
	{{- if .Inspects}}
		unpacked, err := eggtypes.Decode(input)
		if err != nil {
			return err
		}
		env.Logf("middleware: received %#v", unpacked)
		switch input := unpacked.(type) {
		{{- range $inspect := .Inspects}}
		case {{$inspect.GoName}}:
			return m.contract.{{$inspect.GoName}}(
				env,
				{{- range $field := $inspect.Fields}}
					input.{{$field.GoName}},
				{{- end}}
			)
		{{- end}}
		default:
			return fmt.Errorf("middleware: input isn't an inspect")
		}
	{{- else}}
		return fmt.Errorf("inspect not supported")
	{{- end}}
}

// Call eggroll.Roll for the contract using the middleware wrapper.
func Roll(contract iContract) {
	eggroll.Roll(Middleware{contract})
}
`
