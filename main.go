package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"github.com/labstack/echo"
)

const (
	version = "0.1"
)

var (
	h       bool
	v       bool
	cfgpath string
	e       *echo.Echo
	ngx     *Nginx
)

func cmd() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&v, "v", false, "show version")
	flag.StringVar(&cfgpath, "c", "cfg.json", "set cfg and start")
	flag.Parse()

	if v {
		fmt.Println(version)
		os.Exit(0)
	}

	if h {
		flag.Usage()
		os.Exit(0)
	}
}

func InitNgx() {
	ngx = NewNginx(GloabConfig.Bin, GloabConfig.MainConfig,
		GloabConfig.SiteConfigDir, GloabConfig.BackupDir,
		GloabConfig.AllowedPorts)
}

func InitHttp() {
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e = echo.New()
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == GloabConfig.Auth["name"] && password == GloabConfig.Auth["password"] {
			return true, nil
		}
		return false, nil
	}), middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = t
	e.Static("/static", "views")
	e.GET("/", index)
	e.GET("/ping", ping)
	e.GET("/sitecontent", siteContent)
	e.POST("/savesite", saveSite)
	e.GET("/testsite", testSite)
	e.GET("/reloadsite", reloadSite)
	e.GET("/createSite", createSite)
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong from ngxman")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func index(c echo.Context) error {
	sites, err := ngx.SiteConfig()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Sites":   sites,
		"NgxInfo": string(ngx.Info()),
	})
}

func siteContent(c echo.Context) error {
	site := c.QueryParam("site")
	content, err := ngx.SiteConfigContent(site)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"content": string(content),
	})
}

func saveSite(c echo.Context) error {
	site := c.QueryParam("site")
	rename := c.QueryParam("rename")
	content, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	if rename != "" {
		err = ngx.RenameSite(site, rename, content)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	} else {
		err = ngx.SaveSite(site, content)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}
	return c.JSON(http.StatusOK, "success")
}

func testSite(c echo.Context) error {
	out := ngx.TestConfig()
	return c.JSON(http.StatusOK, string(out))
}

func reloadSite(c echo.Context) error {
	out := ngx.Reload()
	return c.JSON(http.StatusOK, string(out))
}

func createSite(c echo.Context) error {
	site := c.QueryParam("site")
	return c.Render(http.StatusOK, "newsite.html", map[string]interface{}{
		"Site": site,
	})
}

func init() {
	// cmd
	cmd()

	// initconfig first
	InitConfig(cfgpath)

	// initngx
	InitNgx()

	// inithttp second
	InitHttp()
}

func main() {
	e.Logger.Fatal(e.Start(GloabConfig.HTTPADDR))
}
