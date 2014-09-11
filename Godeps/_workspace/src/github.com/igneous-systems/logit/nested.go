package logit

import (
	"fmt"
	"strings"
	"sync"
)

type nestedEntry struct {
	pattern   string
	lvl       Level
	expensive *Expensive
}

type NestedData struct {
	lock    sync.Mutex
	inUse   bool
	entries []*nestedEntry
}

func (nd *NestedData) closestEntry(name string) *nestedEntry {
	var best *nestedEntry

	for _, entry := range nd.entries {
		if best != nil && len(best.pattern) > len(entry.pattern) {
			continue
		}
		if len(name) == len(entry.pattern) {
			if name == entry.pattern {
				return entry
			}
		} else if len(name) > len(entry.pattern) {
			if strings.HasPrefix(name, entry.pattern) && name[len(entry.pattern)] == '/' {
				best = entry
			}
		}
	}
	return best
}

func NewNestedData() *NestedData {
	return new(NestedData)
}

func (nd *NestedData) AddEntry(pattern string, lvl Level, expensive *Expensive) error {
	if len(pattern) == 0 {
		return fmt.Errorf("Empty pattern")
	} else if pattern[0] == '/' {
		return fmt.Errorf("Pattern started with illegal character (%s)", pattern)
	} else if pattern[len(pattern)-1] == '/' {
		return fmt.Errorf("Pattern ended with illegal character (%s)", pattern)
	}
	nd.lock.Lock()
	defer nd.lock.Unlock()
	if nd.inUse {
		return fmt.Errorf("Can't modify in use NestedData item")
	}
	for _, entry := range nd.entries {
		if entry.pattern == pattern {
			return fmt.Errorf("Duplicate pattern (%s)", pattern)
		}
	}
	ne := new(nestedEntry)
	ne.pattern = pattern
	ne.lvl = lvl
	nd.entries = append(nd.entries, ne)
	return nil
}

func (nd *NestedData) lockEntries() {
	nd.lock.Lock()
	defer nd.lock.Unlock()
	nd.inUse = true
}
