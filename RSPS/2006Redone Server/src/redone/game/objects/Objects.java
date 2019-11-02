package redone.game.objects;

import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;

@Getter
@Setter
@Accessors(chain = true)
public class Objects {

	public long delay, oDelay;
	public int xp, item, owner, target, times;
	public boolean bait;
	public String belongsTo;
	public int id;
	public int x;
	public int y;
	public int height;
	public int face;
	public int type;
	public int ticks;

	public Objects() {}

	public Objects(int id, int x, int y, int height, int face, int type, int ticks) {
		this.id = id;
		this.x = x;
		this.y = y;
		this.height = height;
		this.face = face;
		this.type = type;
		this.ticks = ticks;
	}
}
