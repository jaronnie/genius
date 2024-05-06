package genius

import (
	"log"
	"testing"
)

func TestEncodeToToml(t *testing.T) {
	source := `
Name = "jaronnie"
`
	g, err := NewFromToml([]byte(source))
	if err != nil {
		log.Fatal(err)
	}

	toml, err := g.EncodeToToml()
	if err != nil {
		log.Fatal(err)
	}
	t.Log(string(toml))
}

func TestEncodeToYaml(t *testing.T) {
	source := `
Name = "jaronnie"
`
	g, err := NewFromToml([]byte(source))
	if err != nil {
		log.Fatal(err)
	}

	yaml, err := g.EncodeToYaml()
	if err != nil {
		log.Fatal(err)
	}
	t.Log(string(yaml))
}

func TestEncodeToJSON(t *testing.T) {
	source := `
Name = "jaronnie"
`
	g, err := NewFromToml([]byte(source))
	if err != nil {
		log.Fatal(err)
	}

	json, err := g.EncodeToJSON()
	if err != nil {
		log.Fatal(err)
	}
	t.Log(string(json))
}
