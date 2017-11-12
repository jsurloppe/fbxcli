package main

import (
	"io/ioutil"
	"log"

	"github.com/jsurloppe/fbxapi"
	"github.com/jsurloppe/fbxcli/cmd"
)

var DEBUG = "NO"

func init() {
	if DEBUG == "NO" {
		log.SetOutput(ioutil.Discard)
		fbxapi.Logr.SetOutput(ioutil.Discard)
	}
}

func main() {
	cmd.Execute()
}
