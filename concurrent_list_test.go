package clist

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func newRandomList(itemsCount int) *ConcurrentList {
	c := New()
	for i := 0; i < itemsCount; i++ {
		c.Add(randSeq(4))
	}

	return c
}

func compareLists(l1, l2 ConcurrentList) bool {
	if len(l1.Items) != len(l2.Items) {
		return false
	}

	for i := 0; i < len(l1.Items); i++ {
		if l1.Items[i] != l2.Items[i] {
			return false
		}
	}

	return true
}

func TestNewCList(t *testing.T) {
	c := New()

	if c == nil {
		t.Errorf("new list instance is null")
	}
	if c.Items == nil {
		t.Errorf("slice not intialized")
	}

	if len(c.Items) != 0 {
		t.Errorf("new list should have length 0, actual %d", len(c.Items))
	}
}

func TestLength(t *testing.T) {
	c := New()
	if c.Length() != 0 {
		t.Errorf("new list should have length of 0, actual %d", c.Length())
	}
	rand.Seed(time.Now().UnixNano())
	itemsCount := rand.Intn(30)
	for i := 0; i < itemsCount; i++ {
		c.Add(randSeq(4))
		if c.Length() != (i + 1) {
			t.Errorf("expecting list's length to be %d, actual %d", (i + 1), c.Length())
		}
	}

	if c.Length() != itemsCount {
		t.Errorf("expecting list's length to be %d, actual %d", itemsCount, c.Length())
	}

	for i := itemsCount - 1; i >= 0; i-- {
		c.Remove(c.Items[i])
		if c.Length() != i {
			t.Errorf("expecting list's length to be %d, actual %d", i, c.Length())
		}
	}

	if c.Length() != 0 {
		t.Errorf("empty list should have length of 0, actual %d", c.Length())
	}
}

func TestContains(t *testing.T) {
	c := New()
	if c.Contains("a") {
		t.Errorf("empty list shouldn't contain any value")
	}

	rand.Seed(time.Now().UnixNano())
	itemsCount := rand.Intn(30)
	items := make([]string, itemsCount)

	// Add some elements to list.
	for i := 0; i < itemsCount; i++ {
		items[i] = randSeq(4)
		c.Add(items[i])
	}

	// Make sure Contains returns true foreach item we've inserted.
	for i := 0; i < itemsCount; i++ {
		if !c.Contains(items[i]) {
			t.Errorf("list should contain item %s", items[i])
		}
	}

	// Remove items.
	for i := 0; i < itemsCount; i++ {
		c.Remove(items[i])
	}

	// Make sure Contains returns false for each item we've removed.
	for i := 0; i < itemsCount; i++ {
		if c.Contains(items[i]) {
			t.Errorf("list shouldn't contain item %s", items[i])
		}
	}
}

func TestAdd(t *testing.T) {
	c := New()

	for i := 9; i >= 0; i-- {
		c.Add(strconv.Itoa(i))
	}
	for idx, item := range c.Items {
		if item != strconv.Itoa(idx) {
			t.Errorf("list items should be sorted, expecting item %s, actual %s", strconv.Itoa(idx), item)
		}
	}
}

func TestRemove(t *testing.T) {
	c := New()

	rand.Seed(time.Now().UnixNano())
	itemsCount := rand.Intn(30)
	items := make([]string, itemsCount)

	// Add some elements to list.
	for i := 0; i < itemsCount; i++ {
		items[i] = randSeq(4)
		c.Add(items[i])
	}

	// Remove items.
	for i := 0; i < itemsCount; i++ {
		c.Remove(items[i])
		if c.Contains(items[i]) {
			t.Errorf("Remove didn't discard of item %s", items[i])
		}
	}

	if len(c.Items) != 0 {
		t.Errorf("Remove didn't discard of all items, found %d items in list, expecting 0", len(c.Items))
	}
}

func TestIndexOf(t *testing.T) {
	c := New()

	if idx := c.indexOf("a"); idx != -1 {
		t.Errorf("index of missing value should be -1, actual %d", idx)
	}

	// Add items 0 to 9.
	for i := 9; i >= 0; i-- {
		c.Add(strconv.Itoa(i))
	}

	for i := 0; i < 10; i++ {
		if c.indexOf(strconv.Itoa(i)) != i {
			t.Errorf("expecting item %s to be at index %d", strconv.Itoa(i), i)
		}
	}

	// Same check as before, only now we've got elements in the list.
	if idx := c.indexOf("a"); idx != -1 {
		t.Errorf("index of missing value should be -1, actual %d", idx)
	}
}

func TestToJSON(t *testing.T) {
	c := New()

	for i := 0; i < 10; i++ {
		c.Add(strconv.Itoa(i))
	}

	expectedJson := "{\"Items\":[\"0\",\"1\",\"2\",\"3\",\"4\",\"5\",\"6\",\"7\",\"8\",\"9\"]}"
	j := c.ToJSON()

	if string(j) != expectedJson {
		t.Errorf("json %s differ from expected %s", string(j), expectedJson)
	}
}

func TestFromJSON(t *testing.T) {
	c := New()
	if c.Contains("a") {
		t.Errorf("empty list shouldn't contain any value")
	}

	rand.Seed(time.Now().UnixNano())
	itemsCount := rand.Intn(30)
	items := make([]string, itemsCount)

	// Add some elements to list.
	for i := 0; i < itemsCount; i++ {
		items[i] = randSeq(4)
		c.Add(items[i])
	}

	j := c.ToJSON()

	cj := FromJSON(j)

	if !compareLists(*c, *cj) {
		t.Errorf("list from json %v differ from origin %v", cj, c)
	}
}

func TestIterBuffered(t *testing.T) {
	c := newRandomList(20)
	counter := 0

	for i := range c.IterBuffered() {
		i = i
		counter++
	}

	if counter != 20 {
		t.Errorf("expecting to scan through 20 items, actual %d", counter)
	}
}
