package genius

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenius(_ *testing.T) {
	source := `
{
	"name":"jaronnie", 
	"age":23, 
	"skills":[
		{
			"Golang":90, 
			"Python":5, 
			"c":5
		}
	]
}
`

	var genius map[string]interface{}

	err := json.Unmarshal([]byte(source), &genius)
	if err != nil {
		fmt.Println(err)
		return
	}

	g := New(genius, WithDelimiter("-"))

	keys := g.GetAllKeys()
	fmt.Println(keys)

	get := g.Get("name")
	fmt.Println(get)

	get = g.Get("skills-0")
	fmt.Println(get)

	get = g.Get("skills-0-Golang")
	fmt.Println(get)

	sub := g.Sub("skills-0")
	fmt.Println(sub.Get("Golang"))

	outerKeys := g.GetOuterKeys("")
	fmt.Println(outerKeys)

	outerKeys = g.GetOuterKeys("skills-0-")
	fmt.Println(outerKeys)

	fmt.Println(g.IsSet("name"))

	_ = g.Set("name", "jaronnie2")

	fmt.Println(g.Get("name"))

	err = g.Set("skills-0-Golang", 100)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(g.Get("skills-0-Golang"))

	err = g.Set("me", []string{"jaron", "gocloudcoder"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(g.Get("me"))
}
