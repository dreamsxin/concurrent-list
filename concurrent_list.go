package clist

import (
	"encoding/json"
	"errors"
	"sync"
)

// ConcurrentList
// Doesn't allow duplicates.

type ConcurrentList struct {
	Items []string
	sync.RWMutex
}

func New() *ConcurrentList {
	return &ConcurrentList{Items: make([]string, 0)}
}

func (c ConcurrentList) Length() int {
	c.Lock()
	defer c.Unlock()
	return len(c.Items)
}

func (c *ConcurrentList) Contains(value string) bool {
	c.RLock()
	defer c.RUnlock()

	// list is empty.
	if len(c.Items) == 0 {
		return false
	}

	return c.indexOf(value) != -1
}

func (c *ConcurrentList) Add(value string) (int, error) {
	c.Lock()
	defer c.Unlock()

	// Incase list is empty.
	if len(c.Items) == 0 {
		c.Items = append(c.Items, value)
		return 0, nil
	}

	// Find insert position.
	var idx int

	// Check for duplicates.
	if idx = c.indexOf(value); idx != -1 {
		return -1, errors.New("value already in list, duplicates are not allowed.")
	}

	// Get insertion index.
	idx = binarySearch(*c, value)

	if c.Items[idx] > value {
		tmp := make([]string, len(c.Items[:idx]))
		copy(tmp, c.Items[:idx])
		tmp = append(tmp, value)
		c.Items = append(tmp, c.Items[idx:]...)
		return idx, nil
	} else {
		tmp := make([]string, len(c.Items[:idx+1]))
		copy(tmp, c.Items[:idx+1])
		tmp = append(tmp, value)
		c.Items = append(tmp, c.Items[idx+1:]...)
		return idx + 1, nil
	}
}

func (c *ConcurrentList) Remove(value string) {
	c.Lock()
	defer c.Unlock()

	idx := c.indexOf(value)
	if idx == -1 {
		return
	}

	c.Items = append(c.Items[:idx], c.Items[idx+1:]...)
}

func (c ConcurrentList) indexOf(value string) int {

	// Empty list.
	if len(c.Items) == 0 {
		return -1
	}

	// Binary search.
	pos := binarySearch(c, value)
	if c.Items[pos] == value {
		return pos
	} else {
		return -1
	}
}

// Returns a buffered iterator which could be used in a for range loop.
func (c ConcurrentList) IterBuffered() <-chan string {
	// Lock.
	c.Lock()
	ch := make(chan string, len(c.Items))
	go func() {
		for _, s := range c.Items {
			ch <- s
		}
		close(ch)
		c.Unlock()
	}()
	return ch
}

func (c ConcurrentList) ToJSON() []byte {
	c.RLock()
	defer c.RUnlock()

	if j, err := json.Marshal(c); err != nil {
		return nil
	} else {
		return j
	}
}

func FromJSON(j []byte) *ConcurrentList {
	c := New()
	if err := json.Unmarshal(j, c); err != nil {
		return nil
	}
	return c
}

func binarySearch(c ConcurrentList, value string) int {
	l := 0                // Left boundry.
	r := len(c.Items) - 1 // Right boundry.

	var i int

	for r > l {
		i = (l + r) / 2
		if c.Items[i] == value {
			return i
		} else if c.Items[i] > value {
			r = i - 1
		} else {
			l = i + 1
		}
	}

	// If we're here that means l == r.
	return l
}
