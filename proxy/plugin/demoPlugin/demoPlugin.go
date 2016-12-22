package demoPlugin

import (
	"github.com/crayontxx/SWProxy-Go/log"
	"github.com/crayontxx/SWProxy-Go/proxy/plugin"
)

type DemoPlugin struct{}

func init() {
	plugin.Register("Demo plugin", DemoPlugin{}, plugin.ReadPlugin, nil)
}

func (p DemoPlugin) OnRequest(m map[string]interface{}) {
	if c, ok := m["command"]; ok {
		log.Printf("REQUEST: %s\n", c)
	}
}

func (p DemoPlugin) OnResponse(m map[string]interface{}) {
	if c, ok := m["command"]; ok {
		log.Printf("RESPONSE: %s\n", c)
	}
}
