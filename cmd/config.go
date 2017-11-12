package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/jsurloppe/fbxapi"
)

const APPID = "com.github.jsurloppe.fbxcli"
const APPNAME = "fbxcli"
const APPVERSION = "0"

// ConfigEntry An entry representing a registered freebox

var ENV struct {
	CfgFile      string
	FreeboxsList map[string]*fbxapi.Freebox
	CurrentAlias string
	Cwd          string
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

func init() {
	clientPool.pool = make(map[string]*fbxapi.Client)
	ENV.Cwd = "/"
}

func NewClientFromPool(alias string) (client *fbxapi.Client, err error) {
	clientPool.mutex.Lock()
	defer clientPool.mutex.Unlock()
	client, ok := clientPool.pool[alias]
	if !ok {
		freebox, ok := ENV.FreeboxsList[alias]
		if !ok {
			return nil, errors.New("Unregistered alias")
		}
		client, err = fbxapi.NewClient(APPID, freebox)
		checkErr(err)
		clientPool.pool[alias] = client
	}
	return client, err
}

func Register(alias string, freebox *fbxapi.Freebox) (client *fbxapi.Client, track_id int, err error) {
	hostname, err := os.Hostname()
	checkErr(err)

	reqAuth := fbxapi.TokenRequest{
		AppId:      APPID,
		AppName:    APPNAME,
		AppVersion: APPVERSION,
		DeviceName: hostname,
	}

	client, err = fbxapi.NewClient(APPID, freebox)
	checkErr(err)
	resp, err := client.Authorize(reqAuth)
	checkErr(err)
	freebox.RespAuthorize = *resp

	// if registered:
	ENV.FreeboxsList[alias] = freebox
	updateConfig()
	track_id = resp.TrackID
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
