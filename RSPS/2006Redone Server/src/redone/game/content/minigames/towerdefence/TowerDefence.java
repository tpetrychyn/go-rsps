package redone.game.content.minigames.towerdefence;

import redone.Server;
import redone.game.content.combat.npcs.NpcEmotes;
import redone.game.npcs.Npc;
import redone.game.npcs.NpcData;
import redone.game.npcs.NpcHandler;
import redone.game.players.Client;
import redone.game.players.Player;
import redone.game.players.PlayerHandler;
import redone.util.Misc;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashSet;
import java.util.Iterator;

class Tile {
    public int x;
    public int y;
    Tile(int x, int y) {
        this.x = x;
        this.y = y;
    }
}

class Hit {
    int delay;
    int damage;
    Hit(int delay, int damage) {
        this.delay = delay;
        this.damage = damage;
    }
}

class Monster {
    public Npc npc;
    public int tile;
    ArrayList<Hit> hits = new ArrayList<>();
    Monster(Npc npc, int tile) {
        this.npc = npc;
        this.tile = tile;
    }
}
public class TowerDefence {
    private int lastSpawn = 0;
    private ArrayList<Tile> tiles = new ArrayList<>();
    private ArrayList<Monster> zombies = new ArrayList<>();
    private ArrayList<Npc> archers = new ArrayList<>();
    private Npc king;
    private int height;
    private Player player;
    private int walkDelay = 0;
    private boolean running = true;

    TowerDefence(Player p) {
        this.player = p;
        this.height = p.heightLevel;
        king = Server.npcHandler.spawnNewNpc(212, 3222, 3225, this.height, 0, 40, 0, 0, 10);
        king.spawnedBy = this.player.playerId;
    }

    void endGame() {
        for (Monster m : zombies) {
            Server.npcHandler.removeNpc(m.npc.npcId);
        }

        for (Npc n : archers) {
            Server.npcHandler.removeNpc(n.npcId);
        }
    }

//    public void NewGame(int x, int y, int height) {
//        tiles.add(0, new Tile(3222, 3217));
//        tiles.add(1, new Tile(3222, 3218));
//        tiles.add(2, new Tile(3222, 3219));
//        tiles.add(3, new Tile(3222, 3220));
//        tiles.add(4, new Tile(3222, 3221));
//    }

    private void spawnZombie(int x, int y) {
        Npc z = Server.npcHandler.spawnNewNpc(73, 3222, 3212, this.height, 0, 18, 0, 0, 10);
        z.spawnedBy = this.player.playerId;
        Monster m = new Monster(z, 0);
        zombies.add(m);
        lastSpawn = 0;
    }

    public void spawnArcher(int x, int y) {
        Npc a = Server.npcHandler.spawnNewNpc(688, x, y, this.height, 0, 100, 3, 10, 1000);
        archers.add(a);
    }

    private void zombieMovement() {
        for (Monster m : new ArrayList<>(zombies)) {
            Npc z = m.npc;
            if (z.isDead) {
                zombies.remove(m);
                continue;
            }
            if (z.absY < 3224 && walkDelay <= 0) {
                z.moveX = NpcHandler.GetMove(z.absX, 3222);
                z.moveY = NpcHandler.GetMove(z.absY, 3224);
                try {
                    z.getNextNPCMovement(z.npcId);
                } catch (Exception e) {
                    this.running = false;
                    return;
                }
                z.updateRequired = true;
                m.tile = m.tile+1;
            } else if (z.absY == 3224) {
                try {
                    int emote = NpcEmotes.getAttackEmote(z.npcId);
                    NpcData.startAnimation(emote, z.npcId);
                } catch (Exception e) {
                    this.running = false;
                    return;
                }
                king.handleHitMask(3);
                king.HP -= 3;
                if (king.HP <= 0) {
                    king.HP = 0;
                }
            }

            for (Hit h : new ArrayList<>(m.hits)) {
                if (h.delay == 1) {
                    z.handleHitMask(h.damage);
                    z.HP -= h.damage;
                    if (z.HP <= 0)  {
                        z.HP = 0;
                        ((Client) player).getItemAssistant().addItem(995, Misc.random(3, 5));
                    }
                    m.hits.remove(h);
                } else {
                    h.delay-=1;
                }
            }
        }
        if (walkDelay <= 0) {
            walkDelay = 1;
        } else {
            walkDelay--;
        }
    }

    public boolean tick() {
        if (lastSpawn >= 5 && zombies.size() <= 5) {
            spawnZombie(0,0);
        } else {
            lastSpawn++;
        }

        zombieMovement();

        for (Npc a : archers) {
            if (a.actionTimer > 0) { continue; }
            a.actionTimer = 3;
            Monster mTarget = findNearest(a.absX, a.absY, 3);
            if (mTarget == null) {
                continue;
            }
            Npc target = mTarget.npc;

            a.turnNpc(target.absX, target.absY);
            NpcData.startAnimation(/*NpcEmotes.getAttackEmote(a.npcId)*/426, a.npcId);
            for (Player player : PlayerHandler.players) {
                if (player != null) {
                    Client c = (Client) player;
                    int nX = a.getX() + NpcHandler.offset(a.npcId);
                    int nY = a.getY() + NpcHandler.offset(a.npcId);
                    int pX = target.getX();
                    int pY = target.getY();
                    int offX = (nY - pY) * -1;
                    int offY = (nX - pX) * -1;
                    c.getActionSender().createProjectile(nX, nY,
                            offX, offY, 50, NpcHandler.getProjectileSpeed(a.npcId),
                            10, 43, 31,
                            0, 65);

                }
            }

            mTarget.hits.add(new Hit(2, Misc.random(1, a.maxHit)));
        }

        if (king.HP <= 0) {
            king.HP = 0;
            this.running = false;
        }

        return this.running;
    }

    private Monster findNearest(int x, int y, int maxDistance) {
        double nearest = Integer.MAX_VALUE;
        Monster target = null;
        for (Monster zombie : zombies) {
            Npc n = zombie.npc;
            if (n.isDead) continue;

            double current = Math.hypot(x - n.absX, y - n.absY);
            if (current > maxDistance) {
                continue;
            }

            if (current < nearest) {
                nearest = current;
                target = zombie;
            }
        }

        return target;
    }
}
