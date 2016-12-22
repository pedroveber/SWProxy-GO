package weakenedUnitPlugin

import "github.com/crayontxx/SWProxy-Go/proxy/plugin"

type WeakenedUnitPlugin struct{}

type AttrCfg struct {
	Attr  string  `json:attr`
	Type  string  `json:type`
	Value float64 `json:value`
}

var cfg struct {
	Tower    []AttrCfg `json:tower`
	Scenario []AttrCfg `json:scenario`
	Dungeon  []AttrCfg `json:dungeon`
	Rift     []AttrCfg `json:rift`
}

func init() {
	plugin.Register("Weakened unit plugin", WeakenedUnitPlugin{}, plugin.WritePlugin, &cfg)
}

func (p WeakenedUnitPlugin) OnRequest(m map[string]interface{}) {
}

func (p WeakenedUnitPlugin) OnResponse(m map[string]interface{}) {
	switch m["command"] {
	case "BattleScenarioStart":
		weakenScenario(m)
	case "BattleDungeonStart":
		weakenDungeon(m)
	case "BattleTrialTowerStart_v2":
		weakenTower(m)
	case "BattleRiftDungeonStart":
		weakenRiftDungeon(m)
	}
}

func weakenScenario(m map[string]interface{}) {
	weaken(m, "opp_unit_list", cfg.Scenario)
}

func weakenDungeon(m map[string]interface{}) {
	weaken(m, "dungeon_unit_list", cfg.Dungeon)
}

func weakenTower(m map[string]interface{}) {
	weaken(m, "trial_tower_unit_list", cfg.Tower)
}

func weakenRiftDungeon(m map[string]interface{}) {
	weaken(m, "rift_dungeon_unit_list", cfg.Rift)
}

func weaken(m map[string]interface{}, listName string, config []AttrCfg) {
	forEachUnit(listName, m, func(unit map[string]interface{}) {
		updateUnit(unit, config)
	})
}

func forEachUnit(listName string, m map[string]interface{}, f func(unit map[string]interface{})) {
	for _, scenario := range m[listName].([]interface{}) {
		for _, unit := range scenario.([]interface{}) {
			f(unit.(map[string]interface{}))
		}
	}
}

func updateUnit(unit map[string]interface{}, config []AttrCfg) {
	for _, c := range config {
		switch c.Type {
		case "percent":
			updateUnitAbilityByPercent(unit, c.Attr, c.Value)
		case "abs":
			updateUnitAbility(unit, c.Attr, c.Value)
		}
	}
}

func updateUnitAbility(unit interface{}, field string, value float64) {
	unit.(map[string]interface{})[field] = value
}

func updateUnitAbilityByPercent(unit interface{}, field string, percent float64) {
	v := unit.(map[string]interface{})[field].(float64)
	unit.(map[string]interface{})[field] = int(v * percent)
}
