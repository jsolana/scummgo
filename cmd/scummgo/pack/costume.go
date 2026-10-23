package pack

import (
	"fmt"
	"path/filepath"
	"strconv"
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

		var act scummgo.CostumeAction
		switch strings.ToLower(anim.Action) {
		case "idle":
			act = scummgo.CostumeIdle(dir())
		case "speak":
			act = scummgo.CostumeSpeak(dir())
		case "walk":
			act = scummgo.CostumeWalk(dir())
		default:
			code, err := strconv.Atoi(anim.Action)
			if err != nil {
				err := fmt.Errorf("neither a default action nor a custom action code: %w", err)
				return fmt.Errorf("invalid action %q: %w", anim.Action, err)
			}
			act = scummgo.CostumeAction(code)
		}
		d.Resource.WithAnimation(act, a)
	}

	return nil
}
