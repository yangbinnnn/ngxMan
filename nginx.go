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
	backupDir     string
	allowedPorts  []string
}

func NewNginx(bin, mainConfig, siteConfigDir, backupDir string, allowedPorts []string) *Nginx {
	return &Nginx{bin, mainConfig, siteConfigDir, backupDir, allowedPorts}
}

func (ngx *Nginx) Backup(site string) {
	p := path.Join(ngx.siteConfigDir, site)
	if !xlib.FileExist(p) {
		return
	}
	bp := path.Join(ngx.backupDir, site+"."+time.Now().Format("200601021504"))
	data, _ := ioutil.ReadFile(p)
	ioutil.WriteFile(bp, data, 0666)
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
		return []byte("加载成功, 访问站点查看")
	}
	return []byte(fmt.Sprintf("%s, %s", "加载失败", out))
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
		if site.IsDir() {
			continue
		}
		if strings.HasPrefix(site.Name(), ".") {
			continue
		}
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

// create or modify site
func (ngx *Nginx) SaveSite(site string, data []byte) error {
	err := ngx.checkConfig(site, data)
	if err != nil {
		return err
	}
	ngx.Backup(site)
	p := path.Join(ngx.siteConfigDir, site)
	return ioutil.WriteFile(p, data, 0666)
}

func (ngx *Nginx) RenameSite(site, newname string, data []byte) error {
	oldpath := path.Join(ngx.siteConfigDir, site)
	newpath := path.Join(ngx.siteConfigDir, newname)
	err := ngx.checkConfig(newname, data)
	if err != nil {
		return err
	}
	ngx.Backup(site)
	err = ioutil.WriteFile(oldpath, data, 0666)
	if err != nil {
		return err
	}
	return os.Rename(oldpath, newpath)
}

func (ngx *Nginx) checkConfig(site string, data []byte) error {
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimSuffix(line, ";")
		if strings.HasPrefix(line, "proxy_pass") {
			if strings.Contains(line, "SERVER_ADDR") {
				return fmt.Errorf("proxy_pass 未正确配置")
			}
		}
		if strings.HasPrefix(line, "server_name") {
			items := strings.Fields(line)
			if len(items) != 2 || items[1] != site {
				return fmt.Errorf("server_name 未正确配置，需要与站点名相同")
			}
		}
		if strings.HasPrefix(line, "listen") {
			items := strings.Fields(line)
			if len(items) != 2 {
				return fmt.Errorf("listen 未正确配置")
			}
			isRightPort := false
			for _, port := range ngx.allowedPorts {
				if items[1] == port {
					isRightPort = true
					break
				}
			}
			if !isRightPort {
				return fmt.Errorf("listen 端口不被允许, 端口限制 %v", ngx.allowedPorts)
			}
		}
	}
	return nil
}
