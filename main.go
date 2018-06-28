package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

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
	ngx = NewNginx(GloabConfig.Bin, GloabConfig.MainConfig, GloabConfig.SiteConfigDir)
}

func InitHttp() {
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e = echo.New()
	e.Renderer = t
	e.Static("/static", "views")
	e.GET("/", index)
	e.GET("/ping", ping)
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
		return c.String(http.StatusInternalServerError, "")
	}
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Sites": sites,
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
