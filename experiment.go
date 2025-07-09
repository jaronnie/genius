package genius

import (
	"strings"

	"github.com/samber/lo"
)

// GetOuterKeys Experiment Feature
func (g *Genius) GetOuterKeys(key string) []string {
	outerKeys := make([]string, 0)
	completeCount := len(strings.Split(key, g.delimiter))
	if completeCount == 0 {
		completeCount = 1
	}

	trimLastDelimiter := strings.TrimRight(key, g.delimiter)
	if g.IsIndexPath(trimLastDelimiter) {
		v := g.Sub(trimLastDelimiter)
		splitKey := strings.Split(trimLastDelimiter, g.delimiter)
		completeCount = completeCount - len(splitKey)
		for _, v := range v.GetAllKeys() {
			split := strings.Split(v, ".")
			if len(split) >= completeCount {
				outerKeys = append(outerKeys, key+strings.Join(split[0:completeCount], g.delimiter))
			}
		}
	} else {
		for _, v := range g.GetAllKeys() {
			split := strings.Split(v, ".")
			if len(split) >= completeCount {
				outerKeys = append(outerKeys, strings.Join(split[0:completeCount], g.delimiter))
			}
		}
	}
	return lo.Uniq(outerKeys)
}
