package ringbuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrowable(t *testing.T) {
	buf := NewGrowable[int]()

	assertBuf := func(hint string, values ...int) {
		t.Helper()
		got := buf.PeekN(999)
		assert.Equal(t, values, got, hint)
	}

	assertDequeue := func(expect int) {
		t.Helper()
		got, valid := buf.Dequeue()
		assert.True(t, valid, "got invalid dequeue")
		assert.Equal(t, expect, got, "invalid dequeue result")
	}

	assertBuf("empty")

	buf.Push(1)
	assertBuf("initial", 1)

	buf.Push(2)
	assertBuf("second", 1, 2)

	buf.Push(3)
	buf.Push(4)
	buf.Push(5)
	buf.Push(6)
	buf.Push(7)
	buf.Push(8)
	assertBuf("after 8", 1, 2, 3, 4, 5, 6, 7, 8)

	buf.Push(9)
	buf.Push(10)
	assertBuf("after growth", 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	assertDequeue(1)
	assertDequeue(2)
	assertBuf("after dequeues", 3, 4, 5, 6, 7, 8, 9, 10)

	buf.Push(11)
	buf.Push(12)
	buf.Push(13)
	buf.Push(14)
	buf.Push(15)
	buf.Push(16)
	buf.Push(17)
	assertBuf("after wrap around", 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17)

	buf.Push(18)
	buf.Push(19)
	assertBuf("after second growth", 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19)

	for i := 3; i <= 19; i++ {
		assertDequeue(i)
	}

	assertBuf("empty again")

	buf.Push(100)
	assertBuf("after starting to fill", 100)
}
