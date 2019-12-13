package main

import (
	"encoding/json"
	"os"
	"log"
	"strings"

	"github.com/btwiuse/pretty"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

type Result struct {
	Id    string
	Value string
}

type Response struct {
	Contents string
}

func main() {
	dec := json.NewDecoder(os.Stdin)
	r := &Result{}
	c := &Response{}
	for dec.More() {
		if err := dec.Decode(r); err != nil {
			panic(err)
		}
		if err := json.Unmarshal([]byte(r.Value), c); err != nil {
			panic(err)
		}
		ast, err := parser.Parse(strings.NewReader(c.Contents))
		if err != nil {
			panic(err)
		}
		if err := Parse(ast.AST); err != nil {
			log.Println(r.Id, err)
			continue
		}
		pretty.Json(r)
	}
}

func Parse(ast *parser.Node) error {
	for _, node := range ast.Children {
		_, err := instructions.ParseInstruction(node)
		if err != nil {
			return err
		}
	}
	return nil
}
