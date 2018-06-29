package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/yangbinnnn/xlib"
)

var SiteNotFound = errors.New("Site not found")
var SiteAlreadyExist = errors.New("Site already exist")

type Nginx struct {
	bin           string
	mainConfig    string
	siteConfigDir string
	allowedPorts  []string
}

func NewNginx(bin, mainConfig, siteConfigDir string, allowedPorts []string) *Nginx {
	return &Nginx{bin, mainConfig, siteConfigDir, allowedPorts}
}

func (ngx *Nginx) Info() (output []byte) {
	cmd := []string{ngx.bin, "-V"}
	_, out := xlib.CmdExec(cmd, 13*time.Second)
	return []byte(out)
}

func (ngx *Nginx) Reload() (output []byte) {
	cmd := []string{ngx.bin, "-s", "reload"}
	ok, out := xlib.CmdExec(cmd, 3*time.Second)
	if ok {
		return []byte("success")
	}
	return []byte(fmt.Sprintf("%s, %s", "failed", out))
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
	err := ngx.checkPort(data)
	if err != nil {
		return err
	}
	p := path.Join(ngx.siteConfigDir, site)
	return ioutil.WriteFile(p, data, 0666)
}

// check site
func (ngx *Nginx) checkPort(data []byte) error {
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "listen") {
			for _, port := range ngx.allowedPorts {
				if strings.Contains(line, port+";") {
					return nil
				}
			}
			return fmt.Errorf("port not allowed, allowd port is %v", ngx.allowedPorts)
		}
	}
	return nil
}
