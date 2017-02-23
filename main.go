package main

import (
	"io/ioutil"
	"log"

	"github.com/jsurloppe/fbxcli/cmd"

	"github.com/jsurloppe/fbxapi"
)

func init() {
	if !cmd.DEBUGMACRO {
		log.SetOutput(ioutil.Discard)
		fbxapi.Logr.SetOutput(ioutil.Discard)
	}
}

func main() {
	cmd.Execute()
}
