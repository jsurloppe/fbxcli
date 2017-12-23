package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/jsurloppe/fbxapi"
)

// ConfigEntry An entry representing a registered freebox

var ENV struct {
	CfgFile      string
	FreeboxsList map[string]*fbxapi.Freebox
	CurrentAlias string
	Cwd          map[string]string
}

var clientPool struct {
	pool  map[string]*fbxapi.Client
	mutex sync.Mutex
}

func PoolLogout() {
	for _, client := range clientPool.pool {
		client.Logout()
	}
}

func getCwd(alias string) string {
	if cwd, ok := ENV.Cwd[alias]; ok {
		return cwd
	}
	return "/"
}

func init() {
	clientPool.pool = make(map[string]*fbxapi.Client)
	ENV.Cwd = make(map[string]string)
}

func NewClient(alias string) (client *fbxapi.Client, err error) {
	clientPool.mutex.Lock()
	defer clientPool.mutex.Unlock()
	client, ok := clientPool.pool[alias]
	if !ok {
		freebox, ok := ENV.FreeboxsList[alias]
		if !ok {
			return nil, errors.New("Unregistered alias")
		}
		client = fbxapi.NewClient(App, freebox)
		clientPool.pool[alias] = client
	}
	return client, err
}

func Register(alias string, freebox *fbxapi.Freebox) (client *fbxapi.Client, track_id int, err error) {
	/*hostname, err := os.Hostname()
	checkErr(err)

	reqAuth := fbxapi.TokenRequest{
		AppId:      APPID,
		AppName:    APPNAME,
		AppVersion: APPVERSION,
		DeviceName: hostname,
	}

	client = fbxapi.NewClient(App, freebox)
	resp, err := client.Authorize(reqAuth)
	checkErr(err)
	freebox.RespAuthorize = *resp

	// if registered:
	ENV.FreeboxsList[alias] = freebox
	updateConfig()
	track_id = resp.TrackID*/
	return
}

func getDefaultFreebox() (alias string) {
	for alias = range ENV.FreeboxsList {
		break
	}
	return
}

func updateConfig() {
	writer, err := os.Create(ENV.CfgFile)
	checkErr(err)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	encoder.Encode(ENV.FreeboxsList)

	writer.Close()
}
