package handler

import (
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet/outgoing"
	"rsps/util"
)

type CombatHandler struct {
	player    *entity.Player
	target    model.Character
	weapon    *util.ItemDefinition
}

func StartCombat(player *entity.Player, target model.Character) {
	combat := &CombatHandler{
		player: player,
	}
	if npc, ok := target.(*entity.Npc); ok {
		combat.target = npc
	}

	player.OngoingAction = combat

	//combat.attack()
	combat.getWeapon()
}

func (c *CombatHandler) Tick() {
	speed := c.weapon.Weapon.AttackSpeed
	if speed == 0 {
		speed = 5
	}
	if c.player.GlobalTickCount == 0 {
	 	c.attack()
	 	c.player.GlobalTickCount = speed
	}
}

func (c *CombatHandler) getWeapon() {
	weaponId := c.player.Equipment.Items[outgoing.EQUIPMENT_SLOTS["weapon"]].ItemId
	weapon := util.GetItemDefinition(weaponId)
	c.weapon = weapon
}

func (c *CombatHandler) attack() {
	c.getWeapon()
	//log.Printf("attacking with %+v", c.weapon)
	c.player.OutgoingQueue = append(c.player.OutgoingQueue, &outgoing.SendSoundPacket{
		Sound:  416,
		Volume: 100,
		Delay:  0,
	})
	c.player.UpdateFlag.SetAnimation(451, 2)

	c.target.GetUpdateFlag().SetAnimation(404, 2)
	c.target.GetUpdateFlag().SetEntityInteraction(c.player)
	c.target.TakeDamage(5)

	if c.target.GetCurrentHitpoints() <= 0 {
		c.player.OngoingAction = nil
		if n, ok := c.target.(*entity.Npc); ok {
			n.Killer = c.player
		}
	}
}
