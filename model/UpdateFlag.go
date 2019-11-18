package model

type UpdateFlag struct {
	UpdateRequired    bool
	NeedsPlacement    bool
	Chat              bool
	ForcedChat        bool
	ForcedMovement    bool
	EntityInteraction bool
	FacePosition      bool
	Appearance        bool

	Animation         bool
	AnimationId       int
	AnimationDuration int

	Graphic           bool
	SingleHit         bool
	DoubleHit         bool
	Transform         bool
}

func (u *UpdateFlag) Clear() {
	u.UpdateRequired = false
	u.NeedsPlacement = false
	u.Chat = false
	u.ForcedChat = false
	u.ForcedMovement = false
	u.EntityInteraction = false
	u.FacePosition = false
	u.Appearance = false
	u.Animation = false
	u.Graphic = false
	u.SingleHit = false
	u.DoubleHit = false
	u.Transform = false

	if u.AnimationDuration >= 0 {
		u.AnimationDuration--
	}
}

func (u *UpdateFlag) SetChat() {
	u.UpdateRequired = true
	u.Chat = true
}

func (u *UpdateFlag) SetAppearance() {
	u.UpdateRequired = true
	u.Appearance = true
}

func (u *UpdateFlag) SetGraphic() {
	u.UpdateRequired = true
	u.Graphic = true
}

func (u *UpdateFlag) SetSingleHit() {
	u.UpdateRequired = true
	u.SingleHit = true
}

func (u *UpdateFlag) SetDoubleHit() {
	u.UpdateRequired = true
	u.DoubleHit = true
}

func (u *UpdateFlag) SetAnimation(id, duration int) {
	u.UpdateRequired = true
	u.Animation = true
	u.AnimationId = id
	u.AnimationDuration = duration
}

func (u *UpdateFlag) ClearAnimation() {
	u.Animation = true
	u.UpdateRequired = true
	u.AnimationId = -1
	u.AnimationDuration = 0
}
