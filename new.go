package genius

import (
	"encoding/json"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

func New(source map[string]interface{}, opts ...Opt) *Genius {
	option := &Option{}
	for _, opt := range opts {
		opt(option)
	}
	defaultOption(option)
	return &Genius{source, option.delimiter}
}

func NewFromRawJSON(source []byte, opts ...Opt) (*Genius, error) {
	var genius map[string]interface{}

	err := json.Unmarshal(source, &genius)
	if err != nil {
		return nil, err
	}
	return New(genius, opts...), nil
}

func NewFromToml(source []byte, opts ...Opt) (*Genius, error) {
	var genius map[string]interface{}

	tree, err := toml.LoadBytes(source)
	if err != nil {
		return nil, err
	}

	genius = tree.ToMap()
	return New(genius, opts...), nil
}

func NewFromYaml(source []byte, opts ...Opt) (*Genius, error) {
	var genius map[string]interface{}

	err := yaml.Unmarshal(source, &genius)
	if err != nil {
		return nil, err
	}

	return New(genius, opts...), nil
}
