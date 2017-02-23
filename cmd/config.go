package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"

	"github.com/jsurloppe/fbxapi"
)

const APPID = "com.github.jsurloppe.fbxcli"
const APPNAME = "fbxcli"
const APPVERSION = "0"

var DEBUGMACRO = false

type CliConfig struct {
	Session string `json:"session,omitempty"`
	Default bool   `json:"default,omitempty"`
	UseSSL  bool   `json:"use_ssl"`
}

// ConfigEntry An entry representing a registered freebox
type ConfigEntry struct {
	fbxapi.Freebox
	fbxapi.RespAuthorize
	CliConfig
}

var ENV struct {
	CfgFile       string
	Freeboxs      map[string]ConfigEntry
	CurrentAlias  string
	CurrentClient *fbxapi.Client
	KeepSession   bool
	Cwd           string
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
	debugStr := os.Getenv("FBXCLI_DEBUG")
	if len(debugStr) > 0 {
		debugBool, err := strconv.ParseBool(debugStr)
		if err == nil {
			DEBUGMACRO = debugBool
		}
	}
	clientPool.pool = make(map[string]*fbxapi.Client)
	ENV.Cwd = "/"
}

func NewClientFromPool(alias string) (client *fbxapi.Client, err error) {
	clientPool.mutex.Lock()
	defer clientPool.mutex.Unlock()
	client, ok := clientPool.pool[alias]
	if !ok {
		freebox, ok := ENV.Freeboxs[alias]
		if !ok {
			return nil, errors.New("Unregistered alias")
		}
		client, err = fbxapi.NewClientFromFreebox(freebox.Freebox, freebox.UseSSL)
		checkErr(err)
		client.SessionToken = freebox.Session
		clientPool.pool[alias] = client
	}
	return client, err
}

func Register(alias string, freebox *fbxapi.Freebox, ssl bool) (client *fbxapi.Client, track_id int, err error) {
	hostname, err := os.Hostname()
	checkErr(err)

	reqAuth := fbxapi.TokenRequest{
		AppId:      APPID,
		AppName:    APPNAME,
		AppVersion: APPVERSION,
		DeviceName: hostname,
	}

	client, err = fbxapi.NewClientFromFreebox(*freebox, ssl)
	checkErr(err)
	resp, err := client.Authorize(reqAuth)
	checkErr(err)
	configEntry := ConfigEntry{Freebox: *freebox, RespAuthorize: *resp}
	// if registered:
	ENV.Freeboxs[alias] = configEntry
	updateConfig()
	track_id = resp.TrackID
	return
}

func updateConfig() {
	writer, err := os.Create(ENV.CfgFile)
	checkErr(err)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	encoder.Encode(ENV.Freeboxs)

	writer.Close()
}
