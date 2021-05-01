package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	path   = "../config/config.json"
	Config = &ConfigData{}
)

type ConfigData struct {
	SocksSwitch      bool          `json:"OpenSocksProxy"`
	MaxPoolConn      int           `json:"MaxPoolConn"`
	ServerPort       string        `json:"ServerPort"`
	ServerHost       string        `json:"ServerHost"`
	ClientPort       string        `json:"ClientPort"`
	Username         string        `json:"Username"`
	Password         string        `json:"Password"`
	IsEncrypted      bool          `json:"IsEncrypted"`
	KeepAlive        bool          `json:"KeepAlive"`
	KeepAliveTimeout time.Duration `json:"KeepAliveTimeout"`
}

func init() {
	path, _ = filepath.Abs(path)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	stream, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	getConfig(stream)
}

func getConfig(s []byte) {
	err := json.Unmarshal(s, Config)
	if err != nil {
		panic(err)
	}
}
