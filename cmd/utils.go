package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jsurloppe/fbxapi"
)

func exitOnPanic() {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			rlshell.writeString(err.Error())
			os.Exit(1)
		}
		panic(r)
	}
}

func recoverOnPanic() {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			rlshell.writeString(err.Error())
		} else {
			panic(r)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getKnownHosts() (knownHosts []string) {
	for _, conf := range ENV.FreeboxsList {
		knownHosts = append(knownHosts, conf.Host)
	}
	return
}

func makePath(current, requested string) string {
	path := strings.TrimSpace(requested)
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("%s/%s", strings.TrimSpace(current), path)
	}
	return path
}

func getCurrentClient() (client *fbxapi.Client, err error) {
	if ENV.FreeboxsList == nil {
		return nil, errors.New("No freebox")
	}
	if len(ENV.CurrentAlias) == 0 {
		for alias := range ENV.FreeboxsList {
			ENV.CurrentAlias = alias
			break
		}
	}
	client, err = NewClientFromPool(ENV.CurrentAlias)
	checkErr(err)
	return client, nil
}
