package model

type UpdateFlag struct {
	UpdateRequired bool
	NeedsPlacement bool
	Chat           bool
	ForcedChat     bool
	ForcedMovement bool

	EntityInteraction bool
	InteractingWith   Character

	Face         bool
	FacePosition *Position

	Appearance bool

	Animation         bool
	AnimationId       int
	AnimationDuration int

	Graphic bool

	SingleHit       bool
	SingleHitDamage int
	DoubleHit       bool
	DoubleHitDamage int

	Transform bool
}

func (u *UpdateFlag) Clear() {
	u.UpdateRequired = false
	u.NeedsPlacement = false
	u.Chat = false
	u.ForcedChat = false
	u.ForcedMovement = false
	u.EntityInteraction = false
	u.Face = false
	u.Appearance = false
	u.Graphic = false

	u.SingleHit = false
	u.SingleHitDamage = 0
	u.DoubleHit = false
	u.DoubleHitDamage = 0

	u.Transform = false
	u.Animation = false

	if u.AnimationDuration >= 0 {
		u.AnimationDuration--
	}
	if u.AnimationDuration <= 0 {
		u.AnimationId = -1
	}
}

func (u *UpdateFlag) SetChat() {
	u.UpdateRequired = true
	u.Chat = true
}

func (u *UpdateFlag) SetEntityInteraction(with Character) {
	u.UpdateRequired = true
	u.InteractingWith = with
	u.EntityInteraction = true
}

func (u *UpdateFlag) SetFacePosition(position *Position) {
	u.FacePosition = position
	u.Face = true
	u.UpdateRequired = true
}

func (u *UpdateFlag) SetAppearance() {
	u.UpdateRequired = true
	u.Appearance = true
}

func (u *UpdateFlag) SetAnimation(id, duration int) {
	u.UpdateRequired = true
	u.Animation = true
	u.AnimationId = id
	u.AnimationDuration = duration
}

func (u *UpdateFlag) ClearAnimation() {
	if u.AnimationId > -1 {
		u.Animation = true
		u.UpdateRequired = true
		u.AnimationId = -1
		u.AnimationDuration = 0
	}
}

func (u *UpdateFlag) SetGraphic() {
	u.UpdateRequired = true
	u.Graphic = true
}

func (u *UpdateFlag) SetSingleHit(damage int) {
	u.UpdateRequired = true
	u.SingleHit = true
	u.SingleHitDamage = damage
}

func (u *UpdateFlag) SetDoubleHit(damage int) {
	u.UpdateRequired = true
	u.DoubleHit = true
	u.DoubleHitDamage = damage
}
