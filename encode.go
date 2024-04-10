package genius

import (
	"github.com/pelletier/go-toml"
)

func (g *Genius) EncodeToToml() ([]byte, error) {
	settings := g.GetAllSettings()

	tree, err := toml.TreeFromMap(settings)
	if err != nil {
		return nil, err
	}
	return tree.Marshal()
}
