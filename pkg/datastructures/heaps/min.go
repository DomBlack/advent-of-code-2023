package heaps

type Item[A any] interface {
	comparable

	// Less returns true if this item is less than the given item.
	Less(A) bool
}

// MinHeap is a min heap of items.
type MinHeap[A Item[A]] struct {
	values  []A            // The heap values
	present map[A]struct{} // A map of values in the heap (we use this to check if a value is in the heap in O(1))
}

// NewMinHeap returns a new min heap.
func NewMinHeap[A Item[A]](capacity ...int) *MinHeap[A] {
	if len(capacity) > 0 {
		return &MinHeap[A]{
			values:  make([]A, 0, capacity[0]),
			present: make(map[A]struct{}, capacity[0]),
		}
	} else {
		return &MinHeap[A]{
			values:  make([]A, 0),
			present: make(map[A]struct{}),
		}
	}
}

func (h *MinHeap[A]) Len() int {
	return len(h.values)
}

// Insert inserts the given item into the heap.
func (h *MinHeap[A]) Insert(item A) {
	h.values = append(h.values, item)
	h.present[item] = struct{}{}
	h.up(len(h.values) - 1)
}

// Remove removes the smallest item from the heap and returns it.
func (h *MinHeap[A]) Remove() A {
	if len(h.values) == 0 {
		panic("cannot remove from empty heap")
	}

	minItem := (h.values)[0]
	(h.values)[0] = (h.values)[len(h.values)-1]
	h.values = (h.values)[:len(h.values)-1]
	h.down(0)

	delete(h.present, minItem)
	return minItem
}

func (h *MinHeap[A]) RemoveItem(item A) {
	if !h.Contains(item) {
		panic("cannot remove item that is not in heap")
	}

	idx := 0
	for i, v := range h.values {
		if v == item {
			idx = i
			break
		}
	}

	(h.values)[idx] = (h.values)[len(h.values)-1]
	h.values = (h.values)[:len(h.values)-1]

	delete(h.present, item)
	h.down(idx)
}

// Update updates the given item in the heap.
func (h *MinHeap[A]) Update(item A) {
	if !h.Contains(item) {
		panic("cannot update item that is not in heap")
	}

	idx := 0
	for i, v := range h.values {
		if v == item {
			idx = i
			break
		}
	}

	h.up(idx)
	h.down(idx)
}

// Contains returns true if the heap contains the given item.
func (h *MinHeap[A]) Contains(item A) bool {
	_, found := h.present[item]
	return found
}

// up moves the item at the given index up the heap.
func (h *MinHeap[A]) up(idx int) {
	for idx > 0 {
		parent := (idx - 1) / 2
		if (h.values)[parent].Less((h.values)[idx]) {
			break
		}

		(h.values)[parent], (h.values)[idx] = (h.values)[idx], (h.values)[parent]
		idx = parent
	}
}

// down moves the item at the given index down the heap.
func (h *MinHeap[A]) down(idx int) {
	for idx < len(h.values) {
		left := idx*2 + 1
		if left >= len(h.values) {
			break
		}

		smallest := left
		if right := left + 1; right < len(h.values) && (h.values)[right].Less((h.values)[left]) {
			smallest = right
		}

		if (h.values)[idx].Less((h.values)[smallest]) {
			break
		}

		(h.values)[idx], (h.values)[smallest] = (h.values)[smallest], (h.values)[idx]
		idx = smallest
	}
}
