package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	//proxy.Verbose = true
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*(com2us.net|qpyou.cn)"))).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Set("X-GoProxy", "yxorPoG-X")
			log.Println(r.URL)
			return r, nil
		})
	log.Fatal(http.ListenAndServe(":443", proxy))
}
