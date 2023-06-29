package genius

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/cast"
)

type Genius struct {
	source    map[string]interface{}
	delimiter string
}

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

func (g *Genius) Get(key string) interface{} {
	return g.get(key)
}

func (g *Genius) GetAllKeys() []string {
	m := map[string]bool{}
	m = g.flattenAndMergeMap(m, g.source, "")
	// convert set of paths to list
	a := make([]string, 0, len(m))
	for x := range m {
		a = append(a, x)
	}
	return a
}

func (g *Genius) GetAllSettings() map[string]interface{} {
	return g.source
}

func (g *Genius) Sub(key string) *Genius {
	data := g.Get(key)
	if data == nil {
		return nil
	}
	subMapx := &Genius{}

	if reflect.TypeOf(data).Kind() == reflect.Map {
		subMapx.source = cast.ToStringMap(data)
		return subMapx
	}
	return nil
}

func (g *Genius) Set(key string, value interface{}) error {
	return g.set(key, value)
}

func (g *Genius) set(key string, value interface{}) error {
	path := strings.Split(key, g.delimiter)
	lastKey := path[len(path)-1]
	configMap := deepSearchStrong(g.source, path[0:len(path)-1])

	switch x := configMap.(type) {
	case map[string]interface{}:
		x[lastKey] = value
	case []interface{}:
		// is this an array
		// lastKey has to be a num
		idx, err := strconv.Atoi(lastKey)
		if err == nil {
			if idx < len(x) && idx >= 0 {
				x[idx] = value
			} else {
				return errors.Errorf("not index [%d] in array", idx)
			}
		}
	default:
		return errors.Errorf("not suppport set key [%v]", key)
	}
	return nil
}

func (g *Genius) IsSet(key string) bool {
	return g.get(key) != nil
}

// Append append data to a slice
func (g *Genius) Append(key string, values ...interface{}) error {
	val := reflect.ValueOf(g.get(key))
	kind := val.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		var sliceValue []interface{}
		length := val.Len()
		if length == 0 {
			return errors.New("not support empty array")
		}
		for j := 0; j < length; j++ {
			sliceValue = append(sliceValue, val.Index(j).Interface())
		}
		for _, v := range values {
			if reflect.TypeOf(v).Kind() != reflect.TypeOf(sliceValue[0]).Kind() {
				return errors.New("not support different typo")
			}
			sliceValue = append(sliceValue, v)
		}
		return g.Set(key, sliceValue)
	}
	return errors.New("only array support append")
}

func (g *Genius) flattenAndMergeMap(shadow map[string]bool, m map[string]interface{}, prefix string) map[string]bool {
	if shadow != nil && prefix != "" && shadow[prefix] {
		// prefix is shadowed => nothing more to flatten
		return shadow
	}
	if shadow == nil {
		shadow = make(map[string]bool)
	}

	var m2 map[string]interface{}
	if prefix != "" {
		prefix += g.delimiter
	}
	for k, val := range m {
		fullKey := prefix + k
		switch v := val.(type) {
		case map[string]interface{}:
			m2 = v
		case map[interface{}]interface{}:
			m2 = cast.ToStringMap(v)
		default:
			// immediate value
			shadow[fullKey] = true
			continue
		}
		// recursively merge to shadow map
		shadow = g.flattenAndMergeMap(shadow, m2, fullKey)
	}
	return shadow
}

func (g *Genius) get(key string) interface{} {
	var val interface{}

	path := strings.Split(key, g.delimiter)
	nested := len(path) > 1

	// Config file next
	val = g.searchIndexableWithPathPrefixes(g.source, path)
	if val != nil {
		return val
	}
	if nested && g.isPathShadowedInDeepMap(path, g.source) != "" {
		return nil
	}

	return nil
}

func (g *Genius) searchIndexableWithPathPrefixes(source interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	// search for path prefixes, starting from the longest one
	for i := len(path); i > 0; i-- {
		prefixKey := strings.Join(path[0:i], g.delimiter)

		var val interface{}
		switch sourceIndexable := source.(type) {
		case []interface{}:
			val = g.searchSliceWithPathPrefixes(sourceIndexable, prefixKey, i, path)
		case map[string]interface{}:
			val = g.searchMapWithPathPrefixes(sourceIndexable, prefixKey, i, path)
		}
		if val != nil {
			return val
		}
	}

	// not found
	return nil
}

// searchSliceWithPathPrefixes searches for a value for path in sourceSlice
//
// This function is part of the searchIndexableWithPathPrefixes recurring search and
// should not be called directly from functions other than searchIndexableWithPathPrefixes.
func (g *Genius) searchSliceWithPathPrefixes(
	sourceSlice []interface{},
	prefixKey string,
	pathIndex int,
	path []string,
) interface{} {
	// if the prefixKey is not a number, or it is out of bounds of the slice
	index, err := strconv.Atoi(prefixKey)
	if err != nil || len(sourceSlice) <= index {
		return nil
	}

	next := sourceSlice[index]

	// Fast path
	if pathIndex == len(path) {
		return next
	}

	switch n := next.(type) {
	case map[interface{}]interface{}:
		return g.searchIndexableWithPathPrefixes(cast.ToStringMap(n), path[pathIndex:])
	case map[string]interface{}, []interface{}:
		return g.searchIndexableWithPathPrefixes(n, path[pathIndex:])
	default:
		// got a value but nested key expected, do nothing and look for next prefix
	}

	// not found
	return nil
}

// searchMapWithPathPrefixes searches for a value for path in sourceMap
//
// This function is part of the searchIndexableWithPathPrefixes recurring search and
// should not be called directly from functions other than searchIndexableWithPathPrefixes.
func (g *Genius) searchMapWithPathPrefixes(
	sourceMap map[string]interface{},
	prefixKey string,
	pathIndex int,
	path []string,
) interface{} {
	next, ok := sourceMap[prefixKey]
	if !ok {
		return nil
	}

	// Fast path
	if pathIndex == len(path) {
		return next
	}

	// Nested case
	switch n := next.(type) {
	case map[interface{}]interface{}:
		return g.searchIndexableWithPathPrefixes(cast.ToStringMap(n), path[pathIndex:])
	case map[string]interface{}, []interface{}:
		return g.searchIndexableWithPathPrefixes(n, path[pathIndex:])
	default:
		// got a value but nested key expected, do nothing and look for next prefix
	}

	// not found
	return nil
}

// isPathShadowedInDeepMap makes sure the given path is not shadowed somewhere
// on its path in the map.
// e.g., if "foo.bar" has a value in the given map, it “shadows”
//
//	"foo.bar.baz" in a lower-priority map
func (g *Genius) isPathShadowedInDeepMap(path []string, m map[string]interface{}) string {
	var parentVal interface{}
	for i := 1; i < len(path); i++ {
		parentVal = g.searchMap(m, path[0:i])
		if parentVal == nil {
			// not found, no need to add more path elements
			return ""
		}
		switch parentVal.(type) {
		case map[interface{}]interface{}:
			continue
		case map[string]interface{}:
			continue
		default:
			// parentVal is a regular value which shadows "path"
			return strings.Join(path[0:i], g.delimiter)
		}
	}
	return ""
}

// searchMap recursively searches for a value for path in source map.
// Returns nil if not found.
// Note: This assumes that the path entries and map keys are lower cased.
func (g *Genius) searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// Fast path
		if len(path) == 1 {
			return next
		}

		// Nested case
		switch v := next.(type) {
		case map[interface{}]interface{}:
			return g.searchMap(cast.ToStringMap(v), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return g.searchMap(v, path[1:])
		default:
			// got a value but nested key expected, return "nil" for not found
			return nil
		}
	}
	return nil
}

func deepSearchStrong(m map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return m
	}
	var currentPath string
	stepArray := false
	var currentArray []interface{}
	var currentEntity interface{}
	for _, k := range path {
		if len(currentPath) == 0 {
			currentPath = k
		} else {
			currentPath = fmt.Sprintf("%v.%v", currentPath, k)
		}
		if stepArray {
			idx, err := strconv.Atoi(k)
			if err != nil {
				return nil
			}
			if len(currentArray) <= idx {
				return nil
			}
			m3, ok := currentArray[idx].(map[string]interface{})
			if !ok {
				return nil
			}
			// continue search from here
			m = m3
			currentEntity = m
			stepArray = false // don't support arrays of arrays
		} else {
			m2, ok := m[k]
			if !ok {
				// intermediate key does not exist
				return nil
			}
			m3, ok := m2.(map[string]interface{})
			if !ok {
				// is this an array
				m4, ok := m2.([]interface{})
				if ok {
					// continue search from here
					currentArray = m4
					currentEntity = currentArray
					stepArray = true
				} else {
					// intermedgiate key is a value
					return nil
				}
			} else {
				// continue search from here
				m = m3
				currentEntity = m
			}
		}
	}

	return currentEntity
}

func (g *Genius) IsIndexPath(key string) bool {
	// defaultKeyDelimiter is .
	split := strings.Split(key, g.delimiter)
	for i, v := range split {
		if _, err := cast.ToIntE(v); err == nil {
			// cast to int successfully
			// Determine whether the key is an array
			join := strings.Join(split[:i], g.delimiter)
			switch g.Get(join).(type) {
			case []interface{}:
				return true
			}
		}
	}
	return false
}
