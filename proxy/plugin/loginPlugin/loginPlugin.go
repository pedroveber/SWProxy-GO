package loginPlugin

import (
	"strings"

	"github.com/crayontxx/SWProxy-Go/proxy/plugin"
)

type LoginPlugin struct{}

func init() {
	plugin.Register("Login plugin", LoginPlugin{}, plugin.WritePlugin, nil)
}

func (p LoginPlugin) OnRequest(m map[string]interface{}) {
}

func (p LoginPlugin) OnResponse(m map[string]interface{}) {
	if _, ok := m["server_url_list"]; !ok {
		return
	}
	list := m["server_url_list"].([]interface{})
	for _, s := range list {
		server := s.(map[string]interface{})
		for k, v := range server {
			switch url := v.(type) {
			case string:
				server[k] = strings.Replace(url, "https", "http", -1)
			}
		}
	}
}
