package main

import (
	"fmt"
	"testing"
)

func TestNginx(t *testing.T) {
	ngx := NewNginx("/usr/local/bin/nginx", "/usr/local/etc/nginx/nginx.conf", "/usr/local/etc/nginx/servers")
	fmt.Println("MainConfig:\n", ngx.MainConfig())
	sites, _ := ngx.SiteConfig()
	fmt.Printf("%v\n", sites)
	fmt.Println("TestConfig:\n", string(ngx.TestConfig()))
	// content, _ := ngx.MainConfigContent()
	// fmt.Println("MainConfigContent:\n", string(content))
	// fmt.Println("Info:\n%s", string(ngx.Info()))
}
