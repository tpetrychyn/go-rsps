package handler

import (
	"rsps/entity"
	"rsps/model"
	"rsps/net/packet/outgoing"
	"rsps/util"
)

type CombatHandler struct {
	attacker model.Character
	target   model.Character
	weapon   *util.ItemDefinition
}

func StartCombat(attacker model.Character, target model.Character) {
	combat := &CombatHandler{
		attacker: attacker,
		target: target,
	}
	attacker.SetOngoingAction(combat)
	if n, ok := target.(*entity.Npc); ok {
		n.Killer = attacker
	}
	combat.getWeapon()
}

func (c *CombatHandler) Tick() {
	speed := c.weapon.Weapon.AttackSpeed
	if speed == 0 {
		speed = 5
	}
	if c.attacker.GetGlobalTickCount() == 0 {
		c.attack()
		c.attacker.SetGlobalTickCount(speed)
	}
}

func (c *CombatHandler) getWeapon() {
	weaponId := 4151//c.attacker.Equipment.Items[outgoing.EQUIPMENT_SLOTS["weapon"]].ItemId
	weapon := util.GetItemDefinition(weaponId)
	c.weapon = weapon
}

func (c *CombatHandler) attack() {
	if c.target.GetCurrentHitpoints() <= 0 {
		c.attacker.SetOngoingAction(nil)

		return
	}

	c.getWeapon()
	//log.Printf("attacking with %+v", c.weapon)
	if p, ok := c.attacker.(*entity.Player); ok {
		p.OutgoingQueue = append(p.OutgoingQueue, &outgoing.SendSoundPacket{
			Sound:  416,
			Volume: 100,
			Delay:  0,
		})
	}

	c.attacker.GetUpdateFlag().SetAnimation(451, 2)

	c.target.GetUpdateFlag().SetAnimation(404, 2)
	c.target.GetUpdateFlag().SetEntityInteraction(c.attacker)
	c.target.TakeDamage(1)

	// TODO: Check auto retaliate, among other issues
	if c.target.GetOngoingAction() == nil {
		StartCombat(c.target, c.attacker)
	}
}
