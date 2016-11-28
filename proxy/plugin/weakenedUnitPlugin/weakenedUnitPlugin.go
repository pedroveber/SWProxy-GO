package weakenedUnitPlugin

import "github.com/crayontxx/SWProxy-Go/proxy/plugin"

type WeakenedUnitPlugin struct{}

func init() {
	plugin.Register("Weakened unit plugin", WeakenedUnitPlugin{}, plugin.WritePlugin)
}

func (p WeakenedUnitPlugin) OnRequest(m map[string]interface{}) {
}

func (p WeakenedUnitPlugin) OnResponse(m map[string]interface{}) {

	if m["command"] == "BattleScenarioStart" {
		for _, scenario := range m["opp_unit_list"].([]interface{}) {
			normalWeakening(scenario)
		}
	} else if m["command"] == "BattleDungeonStart" {
		for _, scenario := range m["dungeon_unit_list"].([]interface{}) {
			normalWeakening(scenario)
		}
	} else if m["command"] == "BattleTrialTowerStart_v2" {
		for _, scenario := range m["trial_tower_unit_list"].([]interface{}) {
			towerWeakening(scenario)
		}
	}
}

func normalWeakening(list interface{}) {
	for _, unit := range list.([]interface{}) {
		// updateUnitAbility(unit, "atk", 11)
		// updateUnitAbility(unit, "def", 11)
		updateUnitAbilityByPercent(unit, "spd", 0.25)
		// updateUnitAbilityByPercent(unit, "con", 0.5)
	}
}

func towerWeakening(list interface{}) {
	for _, unit := range list.([]interface{}) {
		updateUnitAbilityByPercent(unit, "atk", 0.7)
		updateUnitAbilityByPercent(unit, "spd", 0.5)
		updateUnitAbility(unit, "resist", 25)
		updateUnitAbility(unit, "hit_bonus", 10)
		updateUnitAbility(unit, "crit_damage_reduction", 0)

	}
}

func updateUnitAbility(unit interface{}, field string, value float64) {
	unit.(map[string]interface{})[field] = value
}

func updateUnitAbilityByPercent(unit interface{}, field string, percent float64) {
	v := unit.(map[string]interface{})[field].(float64)
	unit.(map[string]interface{})[field] = int(v * percent)
}
