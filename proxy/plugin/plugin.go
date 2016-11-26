package plugin

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/crayontxx/SWProxy-Go/proxy/parser"
)

type Plugin interface {
	OnRequest(m map[string]interface{})
	OnResponse(m map[string]interface{})
}

var (
	readPlugins  = make(map[string]Plugin)
	writePlugins = make(map[string]Plugin)
)

type PluginType int

const (
	ReadPlugin  PluginType = 1
	WritePlugin            = 2
)

func Register(name string, p Plugin, pt PluginType) {
	if pt&ReadPlugin > 0 {
		if _, ok := readPlugins[name]; ok {
			log.Fatalln("Duplicated read plugin name: ", name)
		}
		readPlugins[name] = p
	}
	if pt&WritePlugin > 0 {
		if _, ok := writePlugins[name]; ok {
			log.Fatalln("Duplicated write plugin name: ", name)
		}
		writePlugins[name] = p
	}
}

func getVersion(url *url.URL) int {
	if strings.Contains(url.String(), "_c2.php") {
		return 2
	} else {
		return 1
	}
}

func getRequestPOSTContent(r *http.Request, b []byte) (map[string]interface{}, error) {
	return getPOSTContent(b, parser.DecryptRequest, getVersion(r.URL))
}

func getResponsePOSTContent(r *http.Response, b []byte) (map[string]interface{}, error) {
	return getPOSTContent(b, parser.DecryptResponse, getVersion(r.Request.URL))
}

func getPOSTContent(b []byte, f parser.CryptFunc, version int) (map[string]interface{}, error) {
	b, err := f(b, version)
	if err != nil {
		return nil, err
	}

	var content map[string]interface{}
	err = json.Unmarshal(b, &content)

	return content, err
}

func createRequestPOSTContent(r *http.Request, m map[string]interface{}) ([]byte, error) {
	return createPOSTContent(m, parser.EncryptRequest, getVersion(r.URL))
}

func createResponsePOSTContent(r *http.Response, m map[string]interface{}) ([]byte, error) {
	return createPOSTContent(m, parser.EncryptResponse, getVersion(r.Request.URL))
}

func createPOSTContent(m map[string]interface{}, f parser.CryptFunc, version int) ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return f(b, version)
}

func OnRequest(r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		return
	}

	m, err := getRequestPOSTContent(r, b)
	if err != nil {
		log.Println(err)
		return
	}

	for _, p := range readPlugins {
		p.OnRequest(m)
	}
	for _, p := range writePlugins {
		p.OnRequest(m)
	}

	b, err = createRequestPOSTContent(r, m)
	if err != nil {
		log.Println(err)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	r.ContentLength = int64(len(b))
}

func OnResponse(r *http.Response) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		return
	}

	m, err := getResponsePOSTContent(r, b)
	if err != nil {
		log.Println(err)
		return
	}

	for _, p := range readPlugins {
		p.OnResponse(m)
	}
	for _, p := range writePlugins {
		p.OnResponse(m)
	}

	b, err = createResponsePOSTContent(r, m)
	if err != nil {
		log.Println(err)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	r.ContentLength = int64(len(b))
}