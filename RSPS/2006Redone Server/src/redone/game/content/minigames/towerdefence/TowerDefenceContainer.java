package redone.game.content.minigames.towerdefence;

import redone.game.players.Player;

import java.util.ArrayList;

class TowerDefenceGame {
    Player player;
    int height;
    TowerDefence towerDefence;
}

public class TowerDefenceContainer {
    private ArrayList<TowerDefenceGame> games = new ArrayList<>();

    public TowerDefenceContainer() {

    }

    public TowerDefence getByPlayerId(int pId) {
        for (TowerDefenceGame tdg : games) {
            if (tdg.player.playerId == pId) {
                return tdg.towerDefence;
            }
        }
        return null;
    }

    public void newGame(Player p) {
        TowerDefenceGame tdg = new TowerDefenceGame();
        tdg.player = p;
        tdg.height = p.heightLevel;
        tdg.towerDefence = new TowerDefence(p);
        games.add(tdg);
    }

    public void tick() {
        for (TowerDefenceGame t : new ArrayList<>(games)) {
            boolean status = t.towerDefence.tick();
            if (!status) {
                t.towerDefence.endGame();
                games.remove(t);
            }
        }
    }
}
