package genius

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeniusFromRawJson(t *testing.T) {
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
	g, err := NewFromRawJSON([]byte(source), WithDelimiter("-"))
	assert.Nil(t, err)

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

	err = g.Append("me", "jaronnie", "jaronnie2")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(g.Get("me"))

	err = g.Append("skills", map[string]interface{}{"Golang": "100"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(g.Get("skills"))

	err = g.Set("a", map[string]interface{}{"b": []int{1, 2, 3}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(g.Get("a"))
}

func TestGeniusFromToml(t *testing.T) {
	source := `
Name = "jaronnie"
`
	g, err := NewFromToml([]byte(source))
	assert.Nil(t, err)

	fmt.Println(g.Get("Name"))

	g.Del("Name")

	fmt.Println(g.GetAllSettings())
}
