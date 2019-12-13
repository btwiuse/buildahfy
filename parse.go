package main

import (
	"encoding/json"
	"os"
	"log"
	"fmt"
	"strings"
	"github.com/pkg/errors"

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
		cmds, err := Parse(ast.AST);
		if err != nil {
			panic(err)
		}
		for _, ins := range cmds {
			// log.Println(i)
			// pretty.Json(ins)
			switch c := ins.(type) {
			case *instructions.WorkdirCommand:
				translateWorkdirCommand(c)
			case interface{}:
				continue
			case *instructions.ExposeCommand:
				translateExposeCommand(c)
			case *instructions.StopSignalCommand:
				translateStopSignalCommand(c)
			case *instructions.UserCommand:
				translateUserCommand(c)
			case instructions.Command: // ADD ARG CMD COPY ENTRYPOINT ENV EXPOSE HEALTHCHECK LABEL MAINTAINER ONBUILD RUN SHELL STOPSIGNAL USER VOLUME WORKDIR
				log.Println(c.Name())
			case *instructions.EnvCommand:
				translateEnvCommand(c)
				pretty.JsonString(c)
			case *instructions.Stage: // from command
				if c.Name == "" {
					log.Println("FROM", c.BaseName)
				} else {
					log.Println("FROM", c.BaseName, "as", c.Name)
				}
			default:
				panic(errors.Errorf("%s", "unknown message"))
			}
		}
	}
}

func translateWorkdirCommand(c *instructions.WorkdirCommand) {
	base := "buildah config %s <container>"
	workingdir := fmt.Sprintf("--workingdir %s", c.Path)
	result := fmt.Sprintf(base, workingdir)
	fmt.Println(result)
}

func translateExposeCommand(c *instructions.ExposeCommand) {
	base := "buildah config %s <container>"
	ports := []string{}
	for _, port := range c.Ports {
		port := fmt.Sprintf("--port %s", port)
		ports = append(ports, port)
	}
	result := fmt.Sprintf(base, strings.Join(ports, " "))
	fmt.Println(result)
}

func translateStopSignalCommand(c *instructions.StopSignalCommand) {
	base := "buildah config %s <container>"
	signal := fmt.Sprintf("--stop-signal %s", c.Signal)
	result := fmt.Sprintf(base, signal)
	fmt.Println(result)
}

// TODO: wrap shell variables in single quotes #1250
func translateUserCommand(c *instructions.UserCommand) {
	base := "buildah config %s <container>"
	user := fmt.Sprintf("--user %s", c.User)
	result := fmt.Sprintf(base, user)
	fmt.Println(result)
}

func translateEnvCommand(c *instructions.EnvCommand) {
	base := "buildah config %s <container>"
	envs := []string{}
	for _, kv := range c.Env {
		env := fmt.Sprintf("--env %s=%s", kv.Key, kv.Value)
		envs = append(envs, env)
	}
	result := fmt.Sprintf(base, strings.Join(envs, " "))
	fmt.Println(result)
}

func Parse(ast *parser.Node) ([]interface{}, error) {
	cmds := []interface{}{} // instructions.Command{}
	for _, node := range ast.Children {
		cmd, err := instructions.ParseInstruction(node)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}
