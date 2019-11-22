package handler

import (
	"rsps/entity"
	"time"
)

func Pickpocket(player *entity.Player, npc *entity.Npc) {
	if player.GlobalTickCount == 0 {
		player.GlobalTickCount = 4

		//catchChange := rand.Intn(5)
		if true {
			player.UpdateFlag.SetAnimation(881, 2)
			player.UpdateFlag.SetFacePosition(npc.Position)
			player.GlobalTickCount = 8
			player.IsFrozen = true
			pos := player.Position
			go func() {
				<- time.After(4 * time.Second)
				player.UpdateFlag.SetGraphic(-1)
				player.IsFrozen = false
			}()
			go func() {
				<- time.After(200 * time.Millisecond)
				npc.UpdateFlag.SetAnimation(422, 2)
				npc.UpdateFlag.SetFacePosition(pos)
				npc.UpdateFlag.SetForcedChat("What do you think you're doing?")
				player.UpdateFlag.SetAnimation(404, 2)
				player.UpdateFlag.SetGraphic(80)
			}()
		} else {
			player.UpdateFlag.SetAnimation(881, 2)
			player.UpdateFlag.SetFacePosition(npc.Position)
			player.Inventory.AddItem(995, 3)
		}
	}
}
