package pack

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/jsolana/scummgo"
	"gopkg.in/yaml.v3"
)

// CostumeData is the data for a costume resource.
type CostumeData struct {
	Resource *scummgo.Costume

	workingDir string
}

// NewCostumeData creates a new costume data associated with a working directory.
func NewCostumeData(workingDir string) *CostumeData {
	return &CostumeData{workingDir: workingDir}
}

func (d *CostumeData) UnmarshalYAML(n *yaml.Node) error {
	var data struct {
		Sprites struct {
			Sheet  string
			Width  uint
			Height uint
		}
		Animations []struct {
			Action string
			Dir    string
			Flip   bool
			Frames []struct {
				Row      int
				Columns  []int
				Duration int
			}
		}
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	sprites := scummgo.LoadSpriteSheetFromFile(
		filepath.Join(d.workingDir, data.Sprites.Sheet),
		scummgo.Size{W: int(data.Sprites.Width), H: int(data.Sprites.Height)},
	)
	d.Resource = scummgo.NewCostume(sprites)

	for _, anim := range data.Animations {
		dir := func() scummgo.Direction {
			switch strings.ToLower(anim.Dir) {
			case "right":
				return scummgo.DirRight
			case "left":
				return scummgo.DirLeft
			case "up":
				return scummgo.DirUp
			case "down":
				return scummgo.DirDown
			default:
				panic(fmt.Sprintf("invalid direction %q", anim.Dir))
			}
		}

		a := scummgo.NewAnimation().Flip(anim.Flip)
		for _, frame := range anim.Frames {
			a.AddFrames(
				time.Duration(frame.Duration)*time.Millisecond,
				frame.Row,
				frame.Columns...,
			)
		}

		switch strings.ToLower(anim.Action) {
		case "idle":
			d.Resource.WithAnimation(scummgo.CostumeIdle(dir()), a)
		case "speak":
			d.Resource.WithAnimation(scummgo.CostumeSpeak(dir()), a)
		case "walk":
			d.Resource.WithAnimation(scummgo.CostumeWalk(dir()), a)
		default:
			d.Resource.WithCustomAnimation(anim.Action, a)
		}
	}

	return nil
}
