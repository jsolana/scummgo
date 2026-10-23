package pack

import (
	"path/filepath"

	"github.com/jsolana/scummgo"
	"gopkg.in/yaml.v3"
)

// SoundData is the data for a sound resource.
type SoundData struct {
	Resource *scummgo.SoundTrack

	workingDir string
}

// NewSoundData creates a new sound data associated with a working directory.
func NewSoundData(workingDir string) *SoundData {
	return &SoundData{workingDir: workingDir}
}

func (d *SoundData) UnmarshalYAML(n *yaml.Node) error {
	var data struct {
		Source string
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	d.Resource = scummgo.LoadSoundFromFile(filepath.Join(d.workingDir, data.Source))
	return nil
}
