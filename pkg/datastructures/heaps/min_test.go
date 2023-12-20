package heaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinHeap(t *testing.T) {
	heap := NewMinHeap[testItem]()

	heap.Insert(5)
	heap.Insert(3)
	heap.Insert(7)
	heap.Insert(1)
	heap.Insert(9)
	heap.Insert(2)
	heap.Insert(4)
	heap.Insert(6)
	heap.Insert(8)

	for i := 1; i <= 9; i++ {
		item := heap.Remove()
		assert.Equal(t, i, int(item), "item should be in order")
	}

	heap2 := NewMinHeap[*costItem]()
	item := &costItem{cost: 5}
	heap2.Insert(item)
	heap2.Insert(&costItem{cost: 3})
	heap2.Insert(&costItem{cost: 7})
	item2 := &costItem{cost: 1}
	heap2.Insert(item2)
	heap2.Insert(&costItem{cost: 9})

	item.cost = 2
	heap2.Update(item)
	heap2.Insert(&costItem{cost: 4})

	item2.cost = 6
	heap2.Update(item2)
	heap2.Insert(&costItem{cost: 8})

	lastValue := 0
	for heap2.Len() > 0 {
		item := heap2.Remove()
		assert.True(t, lastValue <= item.cost, "item should be in order")
		lastValue = item.cost

	}
}

type testItem int

func (t testItem) Less(a testItem) bool {
	return t < a
}

type costItem struct {
	cost int
}

func (c *costItem) Less(a *costItem) bool {
	return c.cost < a.cost
}
