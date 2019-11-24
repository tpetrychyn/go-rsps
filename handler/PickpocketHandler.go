package handler

import (
	"log"
	"math/rand"
	"rsps/entity"
	"rsps/model"
	"time"
)

func calculateChance(player *entity.Player, npc *entity.Npc) float32 {
	thievingLevel := player.SkillHelper.Skills[model.Thieving].GetCurrentLevel()
	requiredLevel := 1 // TODO: npc thieving requirements map

	// base chance of 50% success + 2% for every level difference
	chance := 6 + float32(thievingLevel - requiredLevel) / 5
	return chance
}

func Pickpocket(player *entity.Player, npc *entity.Npc) {
	if player.GlobalTickCount == 0 {
		player.GlobalTickCount = 4

		luck := rand.Float32()*(10-0)
		chance := calculateChance(player, npc)
		log.Printf("chance %f, luck %f", chance, luck)
		if true {
			player.UpdateFlag.SetAnimation(881, 2)
			player.UpdateFlag.SetFacePosition(npc.Position)
			player.Inventory.AddItem(995, 3)
			player.SkillHelper.AddExperience(model.Thieving, 100)
		} else {
			player.UpdateFlag.SetAnimation(881, 2)
			player.UpdateFlag.SetFacePosition(npc.Position)
			player.GlobalTickCount = 8
			player.IsFrozen = true
			pos := player.Position
			go func() {
				<- time.After(4 * time.Second)
				// clear freeze after 4 sec
				player.UpdateFlag.SetGraphic(-1)
				player.IsFrozen = false
			}()
			go func() {
				<- time.After(200 * time.Millisecond)
				// hit and freeze the player after 200ms to show the initial attempt
				npc.UpdateFlag.SetAnimation(422, 2)
				npc.UpdateFlag.SetFacePosition(pos)
				npc.UpdateFlag.SetForcedChat("What do you think you're doing?")
				player.UpdateFlag.SetAnimation(404, 2)
				player.UpdateFlag.SetGraphic(80)
			}()
		}
	}
}
