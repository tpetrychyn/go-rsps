package incoming

import (
	"log"
	"rsps/entity"
	"rsps/handler"
	"rsps/model"
	"rsps/net/packet"
)

type ClickNpcPacketHandler struct{}

func (c *ClickNpcPacketHandler) HandlePacket(player *entity.Player, packet *packet.Packet) {
	switch packet.Opcode {
	case ATTACK_NPC_OPCODE:
		c.handleAttackNpc(player, packet)
	case SECOND_CLICK_NPC_OPCODE:
		c.handleNpcSecondClick(player, packet)

	}
}

func (c *ClickNpcPacketHandler) handleAttackNpc(player *entity.Player, packet *packet.Packet) {
	npcIndex := packet.ReadShortA()
	npc := player.GetLoadedNpcs()[int(npcIndex)]
	if npc == nil {
		return
	}
	c.handleAttackNpcInternal(player, npc)
}

func (c *ClickNpcPacketHandler) handleAttackNpcInternal(player *entity.Player, npc model.NpcInterface) {
	if player.DelayedDestination != nil {
		player.DelayedPacket = func() {
			c.handleAttackNpcInternal(player, npc)
		}
		return
	}
	player.UpdateFlag.SetEntityInteraction(npc)
	handler.StartCombat(player, npc)
}

func (c *ClickNpcPacketHandler) handleNpcSecondClick(player *entity.Player, packet *packet.Packet) {
	npcIndex := packet.ReadLEShortA()
	npc := player.GetLoadedNpcs()[int(npcIndex)]
	if npc == nil {
		return
	}
	c.handleNpcSecondClickInternal(player, npc)
}

func (c *ClickNpcPacketHandler) handleNpcSecondClickInternal(player *entity.Player, npc model.NpcInterface) {
	if player.DelayedDestination != nil {
		player.DelayedPacket = func() {
			c.handleNpcSecondClickInternal(player, npc)
		}
		return
	}

	handler.Pickpocket(player, npc.(*entity.Npc))
	log.Printf("second click")
}
