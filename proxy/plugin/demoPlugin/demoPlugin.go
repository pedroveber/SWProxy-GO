package demoPlugin

import (
	"log"

	"github.com/crayontxx/SWProxy-Go/proxy/plugins"
)

type DemoPlugin struct{}

func init() {
	plugin.Register("Demo plugin", DemoPlugin{}, plugin.ReadPlugin)
}

func (p DemoPlugin) OnRequest(m map[string]interface{}) {
	log.Printf("REQUEST: %s\n", m["command"])
}

func (p DemoPlugin) OnResponse(m map[string]interface{}) {
	log.Printf("RESPONSE: %s\n", m["command"])
}
