package redone.world;

import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import com.google.gson.stream.JsonReader;
import redone.Server;
import redone.game.content.skills.core.Mining;
import redone.game.content.skills.core.Woodcutting;
import redone.game.objects.Objects;
import redone.game.players.Client;
import redone.game.players.Player;
import redone.game.players.PlayerHandler;
import redone.util.Misc;
import redone.world.clip.Region;

import java.io.FileNotFoundException;
import java.io.FileReader;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * @author Sanity
 */

public class ObjectHandler {

	private List<Objects> globalObjects = new ArrayList<Objects>();

	public static List<Objects> mapObjects = new ArrayList<Objects>();
	public static List<Objects> removedObjects = new ArrayList<Objects>();

	public void loadWorldObjects() {
        try {
            loadGlobalObjects("C:\\Users\\Taylor\\Desktop\\RSPS\\2006Redone Server\\data\\cfg\\global-objects.json");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
	
	 public Objects getObjectByPosition(int x, int y) {
			for (Objects o : globalObjects) {
                for (Objects globalObject : globalObjects) {
                    if (o.x == x && o.y == y) {
                        return globalObject;
                    }
                }
		}
	    return null;
	 }

	    public void createAnObject(int id, int x, int y, int face) {
	        Objects OBJECT = new Objects(id, x, y, 0, face, 10, 0);
	        if (id == -1) {
	            removeObject(OBJECT);
	        } else {
	            addObject(OBJECT);
	        }
	        //Server.canLoadObjects = true;
	        Server.objectHandler.placeObject(OBJECT);
	    }
		

	public void createAnObject(Client c, int id, int x, int y) {
		Objects OBJECT = new Objects(id, x, y, c.heightLevel, 0, 10, 0);
		if (id == -1) {
			removeObject(OBJECT);
		} else {
			addObject(OBJECT);
		}
		Server.objectHandler.placeObject(OBJECT);
	}

	public void createAnObject(Client c, int id, int x, int y, int face) {
		Objects OBJECT = new Objects(id, x, y, 0, face, 10, 0);
		if (id == -1) {
			removeObject(OBJECT);
		} else {
			addObject(OBJECT);
		}
		Server.objectHandler.placeObject(OBJECT);
	}

	public void createAnObject(int id, int x, int y) {
		Objects OBJECT = new Objects(id, x, y, 0, 0, 10, 0);
		if (id == -1) {
			removeObject(OBJECT);
		} else {
			addObject(OBJECT);
		}
		Server.objectHandler.placeObject(OBJECT);
	}

	/**
	 * Adds object to list
	 **/
	public void addObject(Objects object) {
		globalObjects.add(object);
	}

	/**
	 * Removes object from list
	 **/
	public void removeObject(Objects object) {
		globalObjects.remove(object);
	}

	/**
	 * Does object exist
	 **/
	public Objects objectExists(int objectX, int objectY, int objectHeight) {
		for (Objects o : globalObjects) {
			if (o.getX() == objectX && o.getY() == objectY
					&& o.getHeight() == objectHeight) {
				return o;
			}
		}
		return null;
	}

	/**
	 * Update objects when entering a new region or logging in
	 **/
	public void updateObjects(Client c) {
		for (Objects o : globalObjects) {
			if (c != null) {
				if (c.heightLevel == 0 && o.ticks == 0 && c.distanceToPoint(o.getX(), o.getY()) <= 60) {
					if (Woodcutting.playerTrees(c, o.getId()) || Mining.rockExists(c, o.getId())) {
						c.getActionSender().object(o.getId(), o.getX(), o.getY(), 0, o.getFace(), o.getType());
					}
				}
				if (c.heightLevel == o.getHeight() && !Woodcutting.playerTrees(c, o.getId()) && !Mining.rockExists(c, o.getId()) && o.ticks == 0 && c.distanceToPoint(o.getX(), o.getY()) <= 60) {
					c.getActionSender().object(o.getId(), o.getX(), o.getY(), c.heightLevel, o.getFace(), o.getType());
				}
			}
		}
	}

	/**
	 * Creates the object for anyone who is within 60 squares of the object
	 **/
	public void placeObject(Objects o) {
		for (Player p : PlayerHandler.players) {
			if (p != null) {
				Client person = (Client) p;
				if (person != null) {
					if (person.heightLevel == o.getHeight()
							&& o.ticks == 0) {
						if (person.distanceToPoint(o.getX(),
								o.getY()) <= 60) {
							removeAllObjects(o);
							globalObjects.add(o);
							person.getActionSender().object(
									o.getId(), o.getX(),
									o.getY(), o.getFace(),
									o.getType());
							//Region.addObject(o.getId(), o.getX(), o.getY(), o.getHeight(), o.getType(), o.getFace(), true);
						}
					}
				}
			}
		}
	}

	public void removeAllObjects(Objects o) {
		for (Objects s : globalObjects) {
			if (o.getX() == o.x && o.getY() == o.y
					&& s.getHeight() == o.getHeight()) {
				globalObjects.remove(s);
				break;
			}
		}
	}

	public void process() {
		for (int j = 0; j < globalObjects.size(); j++) {
			if (globalObjects.get(j) != null) {
				Objects o = globalObjects.get(j);
				if (o.ticks > 0) {
					o.ticks--;
				}
				if (o.ticks == 1) {
					Objects deleteObject = objectExists(o.getX(),
							o.getY(), o.getHeight());
					removeObject(deleteObject);
					o.ticks = 0;
					placeObject(o);
					removeObject(o);
					if (isObelisk(o.id)) {
						int index = getObeliskIndex(o.id);
						if (activated[index]) {
							activated[index] = false;
							teleportObelisk(index);
						}
					}
				}
			}

		}
	}

	private void loadGlobalObjects(String fileName) throws FileNotFoundException {
		JsonReader reader = new JsonReader(new FileReader(fileName));
		ArrayList<Objects> objects = new Gson().fromJson(reader, new TypeToken<List<Objects>>(){}.getType());
		globalObjects = objects;

		for (Objects o : objects) {
            Region.addObject(o.id, o.x, o.y, o.height, o.type, o.face, false);
        }

	}

	public final int IN_USE_ID = 14825;

	public boolean isObelisk(int id) {
		for (int obeliskId : obeliskIds) {
			if (obeliskId == id) {
				return true;
			}
		}
		return false;
	}

	public int[] obeliskIds = { 14829, 14830, 111235, 14828, 14826, 14831 };
	public int[][] obeliskCoords = { { 3154, 3618 }, { 3225, 3665 },
			{ 3033, 3730 }, { 3104, 3792 }, { 2978, 3864 }, { 3305, 3914 } };
	public boolean[] activated = { false, false, false, false, false, false };

	public void startObelisk(int obeliskId) {
		int index = getObeliskIndex(obeliskId);
		if (index >= 0) {
			if (!activated[index]) {
				activated[index] = true;
				Objects obby1 = new Objects(14825, obeliskCoords[index][0],
						obeliskCoords[index][1], 0, -1, 10, 0);
				Objects obby2 = new Objects(14825, obeliskCoords[index][0] + 4,
						obeliskCoords[index][1], 0, -1, 10, 0);
				Objects obby3 = new Objects(14825, obeliskCoords[index][0],
						obeliskCoords[index][1] + 4, 0, -1, 10, 0);
				Objects obby4 = new Objects(14825, obeliskCoords[index][0] + 4,
						obeliskCoords[index][1] + 4, 0, -1, 10, 0);
				addObject(obby1);
				addObject(obby2);
				addObject(obby3);
				addObject(obby4);
				Server.objectHandler.placeObject(obby1);
				Server.objectHandler.placeObject(obby2);
				Server.objectHandler.placeObject(obby3);
				Server.objectHandler.placeObject(obby4);
				Objects obby5 = new Objects(obeliskIds[index],
						obeliskCoords[index][0], obeliskCoords[index][1], 0,
						-1, 10, 10);
				Objects obby6 = new Objects(obeliskIds[index],
						obeliskCoords[index][0] + 4, obeliskCoords[index][1],
						0, -1, 10, 10);
				Objects obby7 = new Objects(obeliskIds[index],
						obeliskCoords[index][0], obeliskCoords[index][1] + 4,
						0, -1, 10, 10);
				Objects obby8 = new Objects(obeliskIds[index],
						obeliskCoords[index][0] + 4,
						obeliskCoords[index][1] + 4, 0, -1, 10, 10);
				addObject(obby5);
				addObject(obby6);
				addObject(obby7);
				addObject(obby8);
			}
		}
	}

	public int getObeliskIndex(int id) {
		for (int j = 0; j < obeliskIds.length; j++) {
			if (obeliskIds[j] == id) {
				return j;
			}
		}
		return -1;
	}

	public void teleportObelisk(int port) {
		int random = Misc.random(5);
		while (random == port) {
			random = Misc.random(5);
		}
		for (Player player : PlayerHandler.players) {
			if (player != null) {
				Client c = (Client) player;
				if (Misc.goodDistance(c.getX(), c.getY(),
						obeliskCoords[port][0] + 2, obeliskCoords[port][1] + 2,
						1)) {
					c.getPlayerAssistant().startTeleport(
							obeliskCoords[random][0] + 2,
							obeliskCoords[random][1] + 2, 0, "null");
				}
			}
		}
	}
}
