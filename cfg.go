package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Bin           string            `json:"ngxbin"`
	MainConfig    string            `json:"ngxMainConfig"`
	SiteConfigDir string            `json:"ngxSiteConfigDir"`
	BackupDir     string            `json:"ngxBackupDir"`
	AllowedPorts  []string          `json:"ngxAllowedPorts"`
	HTTPADDR      string            `json:"httpaddr"`
	Auth          map[string]string `json:"auth"`
}

var GloabConfig *Config

func InitConfig(p string) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err.Error())
	}
	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err.Error())
	}
	GloabConfig = config
}
