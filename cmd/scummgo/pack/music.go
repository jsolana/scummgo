package pack

import (
	"path/filepath"

	"github.com/jsolana/scummgo"
	"gopkg.in/yaml.v3"
)

// MusicData is the data for a music resource.
type MusicData struct {
	Resource *scummgo.MusicTrack

	workingDir string
}

// NewMusicData creates a new music data associated with a working directory.
func NewMusicData(workingDir string) *MusicData {
	return &MusicData{workingDir: workingDir}
}

func (d *MusicData) UnmarshalYAML(n *yaml.Node) error {
	var data struct {
		Source string
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	d.Resource = scummgo.LoadMusicFromFile(filepath.Join(d.workingDir, data.Source))
	return nil
}
