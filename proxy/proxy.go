package main

import (
	"flag"
	"net/http"
	"regexp"
	"strconv"

	"github.com/crayontxx/SWProxy-Go/proxy/plugin"

	"github.com/elazarl/goproxy"

	"github.com/crayontxx/SWProxy-Go/log"
	_ "github.com/crayontxx/SWProxy-Go/proxy/plugin/demoPlugin"
	_ "github.com/crayontxx/SWProxy-Go/proxy/plugin/weakenedUnitPlugin"
)

type Options struct {
	port int
}

func (op *Options) parseArguments() {
	flag.IntVar(&op.port, "p", 8080, "port for proxy server")
	flag.Parse()
}

func main() {
	var option Options
	option.parseArguments()

	proxy := goproxy.NewProxyHttpServer()
	com2usUrlRe := regexp.MustCompile("^.*com2us.net/api/gateway_c2.php")
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
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(option.port), proxy))
}
