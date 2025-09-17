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

	err = g.Append("skills", map[string]any{"Golang": "100"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(g.Get("skills"))

	err = g.Set("a", map[string]any{"b": []int{1, 2, 3}})
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

func TestGeniusDelSubstructure(t *testing.T) {
	source := `
{
	"name":"jaronnie", 
	"age":23, 
	"skills":[
		{
			"Golang":90, 
			"Python":5, 
			"c":5
		},
		{
			"Java":80,
			"JavaScript":70
		}
	],
	"profile": {
		"bio": "developer",
		"location": "china"
	}
}
`
	g, err := NewFromRawJSON([]byte(source), WithDelimiter("-"))
	assert.Nil(t, err)

	// Test deleting nested map property
	assert.NotNil(t, g.Get("skills-0-Golang"))
	g.Del("skills-0-Golang")
	assert.Nil(t, g.Get("skills-0-Golang"))
	assert.NotNil(t, g.Get("skills-0-Python")) // Other keys should remain

	// Test deleting array element
	assert.NotNil(t, g.Get("skills-1"))
	g.Del("skills-1")
	assert.Nil(t, g.Get("skills-1"))
	assert.NotNil(t, g.Get("skills-0")) // Remaining element should still exist

	// Test deleting nested object property
	assert.NotNil(t, g.Get("profile-bio"))
	g.Del("profile-bio")
	assert.Nil(t, g.Get("profile-bio"))
	assert.NotNil(t, g.Get("profile-location")) // Other property should remain

	// Test deleting non-existent key (should not panic)
	g.Del("nonexistent-key")
	g.Del("profile-nonexistent")

	fmt.Println("Final state:", g.GetAllSettings())
}

func TestGenius_GetTopLevelKeys(t *testing.T) {
	type fields struct {
		source    map[string]any
		delimiter string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "test1",
			fields: fields{
				source: map[string]any{
					"name": "jaronnie",
					"age":  23,
					"skills": map[string]any{
						"Golang": 90,
						"Python": 5,
						"c":      5,
					},
				},
				delimiter: "-",
			},
			want: []string{"name", "age", "skills"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Genius{
				source:    tt.fields.source,
				delimiter: tt.fields.delimiter,
			}
			assert.ElementsMatch(t, tt.want, g.GetTopLevelKeys(), "GetTopLevelKeys()")
		})
	}
}
