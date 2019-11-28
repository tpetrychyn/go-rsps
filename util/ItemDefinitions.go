package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type WeaponDefinition struct {
	AttackSpeed int `json:"attack_speed"`
	Stances     []struct {
		AttackStyle string      `json:"attack_style"`
		AttackType  string      `json:"attack_type"`
		Boosts      interface{} `json:"boosts"`
		CombatStyle string      `json:"combat_style"`
		Experience  string      `json:"experience"`
	} `json:"stances"`
	WeaponType string `json:"weapon_type"`
}

type EquipmentDefinition struct {
	AttackStab     int         `json:"attack_stab"`
	AttackSlash    int         `json:"attack_slash"`
	AttackCrush    int         `json:"attack_crush"`
	AttackMagic    int         `json:"attack_magic"`
	AttackRanged   int         `json:"attack_ranged"`
	DefenceStab    int         `json:"defence_stab"`
	DefenceSlash   int         `json:"defence_slash"`
	DefenceCrush   int         `json:"defence_crush"`
	DefenceMagic   int         `json:"defence_magic"`
	DefenceRanged  int         `json:"defence_ranged"`
	MeleeStrength  int         `json:"melee_strength"`
	RangedStrength int         `json:"ranged_strength"`
	MagicDamage    int         `json:"magic_damage"`
	Prayer         int         `json:"prayer"`
	Slot           string      `json:"slot"`
	Requirements   interface{} `json:"requirements"`
}

type ItemDefinition struct {
	ID                  int                 `json:"id"`
	Name                string              `json:"name"`
	Incomplete          bool                `json:"incomplete"`
	Members             bool                `json:"members"`
	Tradeable           bool                `json:"tradeable"`
	TradeableOnGe       bool                `json:"tradeable_on_ge"`
	Stackable           bool                `json:"stackable"`
	Noted               bool                `json:"noted"`
	Noteable            bool                `json:"noteable"`
	LinkedIDItem        interface{}         `json:"linked_id_item"`
	LinkedIDNoted       interface{}         `json:"linked_id_noted"`
	LinkedIDPlaceholder int                 `json:"linked_id_placeholder"`
	Placeholder         bool                `json:"placeholder"`
	Equipable           bool                `json:"equipable"`
	EquipableByPlayer   bool                `json:"equipable_by_player"`
	EquipableWeapon     bool                `json:"equipable_weapon"`
	Cost                int                 `json:"cost"`
	Lowalch             int                 `json:"lowalch"`
	Highalch            int                 `json:"highalch"`
	Weight              float64             `json:"weight"`
	BuyLimit            interface{}         `json:"buy_limit"`
	QuestItem           bool                `json:"quest_item"`
	ReleaseDate         string              `json:"release_date"`
	Duplicate           bool                `json:"duplicate"`
	Examine             string              `json:"examine"`
	WikiName            string              `json:"wiki_name"`
	WikiURL             string              `json:"wiki_url"`
	Equipment           EquipmentDefinition `json:"equipment"`
	Weapon              WeaponDefinition    `json:"weapon"`
}

func GetItemDefinition(itemId int) *ItemDefinition {
	return ItemDefinitions[itemId]
}

func GetItemDefinitionByName(name string, noted bool) *ItemDefinition {
	for k, v := range ItemDefinitions {
		if k > 6500 {
			// TODO: seems items above 6??? are "waste disposal"
			continue
		}
		if strings.ToLower(v.Name) == strings.ToLower(name) && v.Noted == noted {
			return v
		}
	}
	return nil
}

var ItemDefinitions map[int]*ItemDefinition

func LoadItemDefinitions() {
	file, err := os.Open("./definitions/items.json")
	if err != nil {
		panic(err)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileBytes, &ItemDefinitions)
	if err != nil {
		panic(err)
	}

	log.Printf("Loaded %d items", len(ItemDefinitions))
}

//
//var max = 7956
//
//func DownloadItems() {
//
//	queue := make(chan ItemDefinition, max)
//	for i := 0; i < 78; i++ {
//		go loadItems(i*102, ((i+1)*102)-1, queue)
//	}
//	items := processItems(queue)
//	sort.Slice(items, func(i, j int) bool {
//		return items[i].ID < items[j].ID
//	})
//
//	file, _ := json.MarshalIndent(items, "", " ")
//
//	_ = ioutil.WriteFile("./definitions/items.json", file, 0644)
//}
//
//func processItems(queue chan ItemDefinition) []ItemDefinition {
//	var items []ItemDefinition
//	for {
//		item := <-queue
//		items = append(items, item)
//		if len(items)%20 == 0 {
//			log.Printf("done %d", len(items))
//		}
//		if len(items) >= 7955 {
//			break
//		}
//	}
//	return items
//}
//
//func loadItems(start, end int, queue chan ItemDefinition) {
//	log.Printf("handling %d to %d", start, end)
//	for i := start; i <= end; i++ {
//		var item ItemDefinition
//		res, _ := resty.New().R().
//			SetResult(&item).
//			Get(fmt.Sprintf("https://www.osrsbox.com/osrsbox-db/items-json/%d.json", i))
//		if res.StatusCode() != 200 {
//			item.BadItem = true
//		}
//		if item.ID == 0 {
//			item.ID = i
//		}
//		queue <- item
//	}
//}
