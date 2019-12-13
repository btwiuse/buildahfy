package main

import (
        "fmt"
        "os"
        "log"
        "strings"
        "encoding/json"

        "github.com/asottile/dockerfile"
        "github.com/btwiuse/pretty"
)

type Result struct {
        Id string
        Value string
}

type Response struct {
        Contents string
}

func main(){
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
                cmds, err := dockerfile.ParseReader(strings.NewReader(c.Contents))
                if err != nil {
                        log.Println(r.Id, err)
                        continue
                }
                pretty.Json(r)
                continue
                // fmt.Println(r.Id)
                for i, cmd := range cmds {
                        if cmd.Cmd == "from" {
                                fmt.Println(i, cmd.Value)
                        }
                }
        }
}

