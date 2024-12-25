package genius

import (
	"encoding/json"
	"fmt"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

func (g *Genius) EncodeToType(configType string) ([]byte, error) {
	settings := g.GetAllSettings()

	switch configType {
	case "json", ".json":
		return json.Marshal(settings)
	case "yaml", "yml", ".yaml", ".yml":
		return yaml.Marshal(settings)
	case "toml", ".toml":
		tree, err := toml.TreeFromMap(settings)
		if err != nil {
			return nil, err
		}
		return tree.Marshal()
	}
	return nil, fmt.Errorf("unsupported config type: %s", configType)
}

func (g *Genius) EncodeToToml() ([]byte, error) {
	settings := g.GetAllSettings()

	tree, err := toml.TreeFromMap(settings)
	if err != nil {
		return nil, err
	}
	return tree.Marshal()
}

func (g *Genius) EncodeToYaml() ([]byte, error) {
	settings := g.GetAllSettings()

	return yaml.Marshal(settings)
}

func (g *Genius) EncodeToJSON() ([]byte, error) {
	settings := g.GetAllSettings()

	return json.Marshal(settings)
}
