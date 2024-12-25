package genius

import (
	"encoding/json"
	"fmt"

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

func NewFromType(source []byte, configType string, opts ...Opt) (*Genius, error) {
	var genius map[string]interface{}

	switch configType {
	case "json", ".json":
		err := json.Unmarshal(source, &genius)
		if err != nil {
			return nil, err
		}
		return New(genius, opts...), nil
	case "yaml", "yml", ".yaml", ".yml":
		err := yaml.Unmarshal(source, &genius)
		if err != nil {
			return nil, err
		}

		return New(genius, opts...), nil
	case "toml", ".toml":
		tree, err := toml.LoadBytes(source)
		if err != nil {
			return nil, err
		}

		genius = tree.ToMap()
		return New(genius, opts...), nil
	}

	return nil, fmt.Errorf("unsupported config type: %s", configType)
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
