package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/yangbinnnn/xlib"
)

var SiteNotFound = errors.New("Site not found")
var SiteAlreadyExist = errors.New("Site already exist")

type Nginx struct {
	bin           string
	mainConfig    string
	siteConfigDir string
}

func NewNginx(bin, mainConfig, siteConfigDir string) *Nginx {
	return &Nginx{bin, mainConfig, siteConfigDir}
}

func (ngx *Nginx) Info() (output []byte) {
	cmd := []string{ngx.bin, "-V"}
	_, out := xlib.CmdExec(cmd, 13*time.Second)
	return []byte(out)
}

func (ngx *Nginx) Reload() (output []byte) {
	cmd := []string{ngx.bin, "-s", "reload"}
	_, out := xlib.CmdExec(cmd, 3*time.Second)
	return []byte(out)
}

func (ngx *Nginx) TestConfig() (output []byte) {
	cmd := []string{ngx.bin, "-t"}
	_, out := xlib.CmdExec(cmd, 3*time.Second)
	return []byte(out)
}

// nginx.conf
func (ngx *Nginx) MainConfig() string {
	return ngx.mainConfig
}

func (ngx *Nginx) MainConfigContent() ([]byte, error) {
	return ioutil.ReadFile(ngx.mainConfig)
}

// conf.d/site.domain.com
func (ngx *Nginx) SiteConfig() ([]string, error) {
	files, err := ioutil.ReadDir(ngx.siteConfigDir)
	if err != nil {
		return nil, err
	}
	var sites []string
	for _, site := range files {
		sites = append(sites, site.Name())
	}
	return sites, nil
}

func (ngx *Nginx) SiteConfigContent(site string) ([]byte, error) {
	p := path.Join(ngx.siteConfigDir, site)
	if !xlib.FileExist(p) {
		return nil, SiteNotFound
	}
	return ioutil.ReadFile(p)
}

func (ngx *Nginx) DeleteSite(site string) error {
	p := path.Join(ngx.siteConfigDir, site)
	if !xlib.FileExist(p) {
		return nil
	}
	return os.Remove(p)
}

// create and modify site
func (ngx *Nginx) SaveSite(site string, data []byte) error {
	p := path.Join(ngx.siteConfigDir, site)
	return ioutil.WriteFile(p, data, 0666)
}
