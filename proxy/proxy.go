package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/crayontxx/SWProxy-Go/proxy/plugin"

	"github.com/elazarl/goproxy"

	"github.com/crayontxx/SWProxy-Go/log"
	_ "github.com/crayontxx/SWProxy-Go/proxy/plugin/demoPlugin"
	_ "github.com/crayontxx/SWProxy-Go/proxy/plugin/loginPlugin"
	_ "github.com/crayontxx/SWProxy-Go/proxy/plugin/weakenedUnitPlugin"
)

var cfg struct {
	Port    int
	Plugins json.RawMessage
}

func fatalOr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// TODO: avoid hardcoding
	f, err := os.Open("./config.cfg")
	fatalOr(err)
	bs, err := ioutil.ReadAll(f)
	fatalOr(err)
	err = json.Unmarshal(bs, &cfg)
	fatalOr(err)
	plugin.ApplyConfig(cfg.Plugins)

	proxy := goproxy.NewProxyHttpServer()
	com2usUrlRe := regexp.MustCompile("^.*((com2us.net|qpyou.cn)/api/gateway_c2.php|location_c2.php)")
	proxy.OnRequest(goproxy.UrlMatches(com2usUrlRe)).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			plugin.OnRequest(r)
			return r, nil
		})

	proxy.OnResponse(goproxy.UrlMatches(com2usUrlRe)).DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			plugin.OnResponse(resp)
			return resp
		})
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), proxy))
}
