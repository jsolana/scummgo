package scummgo

import (
	"io"
	"sort"
)

// CostumeAction is a value that represents an action for a costume. For predefined actions idle,
// speak, and walk, use the CostumeIdle, CostumeSpeak, and CostumeWalk functions respectively to
// refer to them. Custom animations are identified by name instead, see WithCustomAnimation.
type CostumeAction byte

// CostumeIdle returns a costume action for the idle action in the given direction.
func CostumeIdle(dir Direction) CostumeAction {
	return CostumeAction((0 << 2) | (dir & 0x03))
}

// CostumeSpeak returns a costume action for the speak action in the given direction.
func CostumeSpeak(dir Direction) CostumeAction {
	return CostumeAction((1 << 2) | (dir & 0x03))
}

// CostumeWalk returns a costume action for the walk action in the given direction.
func CostumeWalk(dir Direction) CostumeAction {
	return CostumeAction((2 << 2) | (dir & 0x03))
}

// Costume is a struct that represents a costume for an actor or a room animation.
type Costume struct {
	sprites *SpriteSheet

	anims  map[CostumeAction]*Animation
	custom map[string]*Animation
}

// NewCostume creates a new costume.
func NewCostume(sprites *SpriteSheet) *Costume {
	return &Costume{
		sprites: sprites,
		anims:   make(map[CostumeAction]*Animation),
		custom:  make(map[string]*Animation),
	}
}

// WithAnimationStand sets the stand animation for the given direction.
func (c *Costume) WithAnimation(act CostumeAction, anim *Animation) *Costume {
	c.anims[act] = anim
	return c
}

// WithCustomAnimation sets a custom animation identified by name.
func (c *Costume) WithCustomAnimation(name string, anim *Animation) *Costume {
	c.custom[name] = anim
	return c
}

// BinaryEncode encodes the costume to a binary format. The format is as follows:
// - sprite sheet.
// - uint32: the number of animations.
// - for each animation:
//   - byte: the action.
//   - the animation.
//
// - uint32: the number of custom animations.
// - for each custom animation:
//   - string: the name.
//   - the animation.
func (c *Costume) BinaryEncode(w io.Writer) (n int, err error) {
	n, err = BinaryEncode(w, c.sprites, uint32(len(c.anims)))
	if err != nil {
		return n, err
	}
	// Sort actions for deterministic on-disk layout.
	actions := make([]byte, 0, len(c.anims))
	for act := range c.anims {
		actions = append(actions, byte(act))
	}
	sort.Slice(actions, func(i, j int) bool { return actions[i] < actions[j] })
	for _, act := range actions {
		nn, err := BinaryEncode(w, act, c.anims[CostumeAction(act)])
		n += nn
		if err != nil {
			return n, err
		}
	}
	nn, err := BinaryEncode(w, uint32(len(c.custom)))
	n += nn
	if err != nil {
		return n, err
	}
	names := make([]string, 0, len(c.custom))
	for name := range c.custom {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		nn, err := BinaryEncode(w, name, c.custom[name])
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// BinaryDecode decodes the costume from a binary format. See BinaryEncode for the format.
func (c *Costume) BinaryDecode(r io.Reader) error {
	c.sprites = new(SpriteSheet)
	c.anims = make(map[CostumeAction]*Animation)
	c.custom = make(map[string]*Animation)

	var count uint32
	if err := BinaryDecode(r, c.sprites, &count); err != nil {
		return err
	}
	for i := uint32(0); i < count; i++ {
		var act byte
		anim := new(Animation)
		if err := BinaryDecode(r, &act, anim); err != nil {
			return err
		}
		c.anims[CostumeAction(act)] = anim
	}
	if err := BinaryDecode(r, &count); err != nil {
		return err
	}
	for i := uint32(0); i < count; i++ {
		var name string
		anim := new(Animation)
		if err := BinaryDecode(r, &name, anim); err != nil {
			return err
		}
		c.custom[name] = anim
	}
	return nil
}

func (c *Costume) draw(act CostumeAction, pos Position, scale float32) {
	if anim := c.anims[act]; anim != nil {
		anim.Draw(c.sprites, pos, scale)
	}
}
