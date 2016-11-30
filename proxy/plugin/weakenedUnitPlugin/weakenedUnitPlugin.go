package weakenedUnitPlugin

import (
	"log"

	"github.com/crayontxx/SWProxy-Go/proxy/plugin"
)

type WeakenedUnitPlugin struct{}

func init() {
	plugin.Register("Weakened unit plugin", WeakenedUnitPlugin{}, plugin.WritePlugin)
}

func (p WeakenedUnitPlugin) OnRequest(m map[string]interface{}) {
	switch m["command"] {
	case "BattleRiftDungeonResult":
		log.Println(m)
	}
}

func (p WeakenedUnitPlugin) OnResponse(m map[string]interface{}) {
	switch m["command"] {
	case "BattleScenarioStart":
		normalWeakening(m, "opp_unit_list")
	case "BattleDungeonStart":
		normalWeakening(m, "dungeon_unit_list")
	case "BattleTrialTowerStart_v2":
		towerWeakening(m, "trial_tower_unit_list")
	case "BattleRiftDungeonStart":
		riftDungeonWeakening(m, "rift_dungeon_unit_list")
	case "BattleRiftDungeonResult":
		log.Println(m)
	}
}

func normalWeakening(m map[string]interface{}, listName string) {
	forEachUnit(listName, m, func(unit map[string]interface{}) {
		updateUnitAbilityByPercent(unit, "spd", 0.5)
		updateUnitAbilityByPercent(unit, "con", 0.2)
		updateUnitAbilityByPercent(unit, "def", 0.01)
		// updateUnitAbilityByPercent(unit, "def", 0.1)
	})
}

func towerWeakening(m map[string]interface{}, listName string) {
	forEachUnit(listName, m, func(unit map[string]interface{}) {
		updateUnitAbility(unit, "atk", 100)
		updateUnitAbilityByPercent(unit, "spd", 0.5)
		updateUnitAbility(unit, "resist", 25)
		updateUnitAbility(unit, "hit_bonus", 10)
		updateUnitAbility(unit, "crit_damage_reduction", 0)
	})
}

func riftDungeonWeakening(m map[string]interface{}, listName string) {
	forEachUnit(listName, m, func(unit map[string]interface{}) {
		updateUnitAbilityByPercent(unit, "atk", 0.5)
		updateUnitAbilityByPercent(unit, "spd", 0.2)
		updateUnitAbility(unit, "resist", 15)
	})
}

func forEachUnit(listName string, m map[string]interface{}, f func(unit map[string]interface{})) {
	for _, scenario := range m[listName].([]interface{}) {
		for _, unit := range scenario.([]interface{}) {
			f(unit.(map[string]interface{}))
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
