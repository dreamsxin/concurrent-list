# concurrent list [![Circle CI](https://circleci.com/gh/streamrail/concurrent-list.svg?style=svg)](https://circleci.com/gh/streamrail/concurrent-list)


A concurrent safe slice of type string, the list doesn't allow duplicates.
Items within the list are sorted, so don't expect to find your Nth item to be at index N-1, see example below.

## usage

Import the package:

```go
import (
	"github.com/streamrail/concurrent-list"
)

```

## example


```go

	c := clist.New()
	c.Add("foo")
	c.Add("bar")

	// Prints "bar", "foo", list doesn't necessarily maintains insert order.
	for item := range c.IterBuffered() {
		fmt.Println(item)
	}	
	
	j := c.ToJSON()
	dup := clist.FromJSON(j)


	if c.Contains("foo") {
		c.Remove("foo")
	}

	if c.Length() > 0 {
		c.Remove("bar")
	}

```

For more examples have a look at concurrent_list_test.go


Running tests:
```bash
go test "github.com/streamrail/concurrent-list"
```