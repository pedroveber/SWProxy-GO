package plugin

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/crayontxx/SWProxy-Go/log"
	"github.com/crayontxx/SWProxy-Go/proxy/parser"
)

type Plugin interface {
	OnRequest(m map[string]interface{})
	OnResponse(m map[string]interface{})
}

type PluginPortfolio struct {
	P      Plugin
	Type   PluginType
	Config interface{} // ptr
}

type PluginConfig struct {
	Name string
	Data json.RawMessage
}

var (
	plugins = make(map[string]*PluginPortfolio)
)

type PluginType int

const (
	ReadPlugin  PluginType = 1
	WritePlugin            = 2
)

func Register(name string, p Plugin, pt PluginType, config interface{}) {
	if _, ok := plugins[name]; ok {
		log.Fatalln("Duplicated plguin name: ", name)
	}
	value := reflect.ValueOf(config)
	if config != nil && value.Kind() != reflect.Ptr {
		log.Fatalf("Error in plugin %v: need a pointer to config variable", name)
	}
	plugins[name] = &PluginPortfolio{p, pt, config}
	log.Println("Loaded Plugin: ", name)
}

func ApplyConfig(cfg json.RawMessage) {
	var pcs []PluginConfig
	err := json.Unmarshal(cfg, &pcs)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range pcs {
		if _, ok := plugins[c.Name]; !ok {
			log.Debug("skip config of plugin: ", c.Name)
			continue
		}
		pf := plugins[c.Name]
		err = json.Unmarshal(c.Data, pf.Config)
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Loaded config of plugin(%v) %+v\n", c.Name, pf.Config)
	}
}

func getVersion(url *url.URL) int {
	if strings.Contains(url.String(), "_c2.php") {
		return 2
	} else {
		return 1
	}
}

func getRequestPOSTContent(r *http.Request, b []byte) (map[string]interface{}, []byte, error) {
	return getPOSTContent(b, parser.DecryptRequest, getVersion(r.URL))
}

func getResponsePOSTContent(r *http.Response, b []byte) (map[string]interface{}, []byte, error) {
	return getPOSTContent(b, parser.DecryptResponse, getVersion(r.Request.URL))
}

func getPOSTContent(b []byte, f parser.CryptFunc, version int) (map[string]interface{}, []byte, error) {
	b, err := f(b, version)
	if err != nil {
		return nil, nil, err
	}

	var content map[string]interface{}
	err = json.Unmarshal(b, &content)

	return content, b, err
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
	if r.Method != "POST" {
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Debugln(err)
		return
	}

	m, b, err := getRequestPOSTContent(r, b)
	if err != nil {
		log.Debugln(err)
		return
	}

	log.Debugln("Request:", string(b))
	for _, pf := range plugins {
		if pf.Type == ReadPlugin {
			pf.P.OnRequest(m)
		}
	}

	for _, pf := range plugins {
		if pf.Type == WritePlugin {
			pf.P.OnRequest(m)
		}
	}

	b, err = createRequestPOSTContent(r, m)
	if err != nil {
		log.Debugln(err)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	r.ContentLength = int64(len(b))
}

func OnResponse(r *http.Response) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Debugln(err)
		return
	}

	m, b, err := getResponsePOSTContent(r, b)
	if err != nil {
		log.Debugln(err)
		return
	}

	log.Debugln("Response:", string(b))
	for _, pf := range plugins {
		if pf.Type == ReadPlugin {
			pf.P.OnResponse(m)
		}
	}

	for _, pf := range plugins {
		if pf.Type == WritePlugin {
			pf.P.OnResponse(m)
		}
	}

	b, err = createResponsePOSTContent(r, m)
	if err != nil {
		log.Debugln(err)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	r.ContentLength = int64(len(b))
}
