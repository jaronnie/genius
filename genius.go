package genius

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type Genius struct {
	source    map[string]any
	delimiter string
}

func (g *Genius) Get(key string) any {
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

func (g *Genius) GetTopLevelKeys() []string {
	a := make([]string, 0, len(g.source))
	for k := range g.source {
		a = append(a, k)
	}
	return a
}

func (g *Genius) GetAllSettings() map[string]any {
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

func (g *Genius) Del(key string) {
	path := strings.Split(key, g.delimiter)

	// Handle simple case for root level keys
	if len(path) == 1 {
		delete(g.source, key)
		return
	}

	// Handle nested keys
	lastKey := path[len(path)-1]
	parentPath := path[0 : len(path)-1]

	// Find the parent container
	parentContainer := deepSearchStrong(g.source, parentPath)
	if parentContainer == nil {
		// Parent doesn't exist, nothing to delete
		return
	}

	switch x := parentContainer.(type) {
	case map[string]any:
		delete(x, lastKey)
	case []any:
		// Handle array deletion by index
		idx, err := strconv.Atoi(lastKey)
		if err != nil || idx < 0 || idx >= len(x) {
			// Invalid index, nothing to delete
			return
		}
		// Remove element at index by creating a new slice without that element
		newSlice := make([]any, 0, len(x)-1)
		newSlice = append(newSlice, x[:idx]...)
		newSlice = append(newSlice, x[idx+1:]...)

		// Update the parent with the new slice
		if len(parentPath) == 0 {
			// This shouldn't happen as we handle root case above, but just in case
			return
		}
		grandParentPath := parentPath[0 : len(parentPath)-1]
		grandParentKey := parentPath[len(parentPath)-1]
		grandParent := deepSearchStrong(g.source, grandParentPath)

		if grandParentMap, ok := grandParent.(map[string]any); ok {
			grandParentMap[grandParentKey] = newSlice
		}
	}
}

func (g *Genius) Set(key string, value any) error {
	return g.set(key, value)
}

func (g *Genius) set(key string, value any) error {
	path := strings.Split(key, g.delimiter)
	lastKey := path[len(path)-1]
	configMap := deepSearchStrong(g.source, path[0:len(path)-1])

	if configMap == nil {
		configMap = deepSearch(g.source, path[0:len(path)-1])
	}

	switch x := configMap.(type) {
	case map[string]any:
		x[lastKey] = value
	case []any:
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
func (g *Genius) Append(key string, values ...any) error {
	val := reflect.ValueOf(g.get(key))
	kind := val.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		var sliceValue []any
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

func (g *Genius) flattenAndMergeMap(shadow map[string]bool, m map[string]any, prefix string) map[string]bool {
	if shadow != nil && prefix != "" && shadow[prefix] {
		// prefix is shadowed => nothing more to flatten
		return shadow
	}
	if shadow == nil {
		shadow = make(map[string]bool)
	}

	var m2 map[string]any
	if prefix != "" {
		prefix += g.delimiter
	}
	for k, val := range m {
		fullKey := prefix + k
		switch v := val.(type) {
		case map[string]any:
			m2 = v
		case map[any]any:
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

func (g *Genius) get(key string) any {
	var val any

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

func (g *Genius) searchIndexableWithPathPrefixes(source any, path []string) any {
	if len(path) == 0 {
		return source
	}

	// search for path prefixes, starting from the longest one
	for i := len(path); i > 0; i-- {
		prefixKey := strings.Join(path[0:i], g.delimiter)

		var val any
		switch sourceIndexable := source.(type) {
		case []any:
			val = g.searchSliceWithPathPrefixes(sourceIndexable, prefixKey, i, path)
		case map[string]any:
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
	sourceSlice []any,
	prefixKey string,
	pathIndex int,
	path []string,
) any {
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
	case map[any]any:
		return g.searchIndexableWithPathPrefixes(cast.ToStringMap(n), path[pathIndex:])
	case map[string]any, []any:
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
	sourceMap map[string]any,
	prefixKey string,
	pathIndex int,
	path []string,
) any {
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
	case map[any]any:
		return g.searchIndexableWithPathPrefixes(cast.ToStringMap(n), path[pathIndex:])
	case map[string]any, []any:
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
func (g *Genius) isPathShadowedInDeepMap(path []string, m map[string]any) string {
	var parentVal any
	for i := 1; i < len(path); i++ {
		parentVal = g.searchMap(m, path[0:i])
		if parentVal == nil {
			// not found, no need to add more path elements
			return ""
		}
		switch parentVal.(type) {
		case map[any]any:
			continue
		case map[string]any:
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
func (g *Genius) searchMap(source map[string]any, path []string) any {
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
		case map[any]any:
			return g.searchMap(cast.ToStringMap(v), path[1:])
		case map[string]any:
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

// deepSearch scans deep maps, following the key indexes listed in the
// sequence "path".
// The last value is expected to be another map, and is returned.
//
// In case intermediate keys do not exist, or map to a non-map value,
// a new map is created and inserted, and the search continues from there:
// the initial map "m" may be modified!
func deepSearch(m map[string]any, path []string) map[string]any {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]any)
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]any)
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]any)
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}

func deepSearchStrong(m map[string]any, path []string) any {
	if len(path) == 0 {
		return m
	}
	var currentPath string
	stepArray := false
	var currentArray []any
	var currentEntity any
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
			m3, ok := currentArray[idx].(map[string]any)
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
			m3, ok := m2.(map[string]any)
			if !ok {
				// is this an array
				m4, ok := m2.([]any)
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
			case []any:
				return true
			}
		}
	}
	return false
}
