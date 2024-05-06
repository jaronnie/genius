package genius

import (
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

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
