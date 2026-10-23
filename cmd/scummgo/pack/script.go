package pack

import (
	"fmt"

	"github.com/jsolana/scummgo"
	"gopkg.in/yaml.v3"
)

type ScriptLanguage string

const (
	ScriptLanguageLua ScriptLanguage = "lua"
)

// ScriptData is the data for a script resource.
type ScriptData struct {
	Resource *scummgo.Script
}

func (d *ScriptData) UnmarshalYAML(n *yaml.Node) error {
	var data struct {
		Language string
		Code     string
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	switch data.Language {
	case "", "lua", "Lua", "LUA":
		d.Resource = scummgo.NewScript(scummgo.ScriptLua, []byte(data.Code))
	default:
		return fmt.Errorf("unknown script language: %s", data.Language)
	}

	return nil
}
