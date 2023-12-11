package day05

import (
	"fmt"
	"sort"

	"github.com/cockroachdb/errors"
)

type Offsets []Offset

// FindIdx returns the index of the offset which contains the given src position
//
// If no offset contains the given src position then -1 is returned
func (o Offsets) FindIdx(src int) (idx int) {
	if len(o) == 0 {
		return -1
	}

	idx = sort.Search(len(o), func(i int) bool {
		return o[i].InputTo >= src
	})

	if idx < len(o) {
		return idx
	} else {
		return -1
	}
}

// Compress merges any offsets next to each other which have the same offset
func (o Offsets) Compress() Offsets {
	if len(o) == 0 {
		return nil
	}

	rtn := make(Offsets, 0, len(o))

	lastIdx := 0
	rtn = append(rtn, o[0])

	for i := 1; i < len(o); i++ {
		this := o[i]

		if this.OffsetBy == rtn[lastIdx].OffsetBy {
			rtn[lastIdx] = this
		} else {
			rtn = append(rtn, this)
			lastIdx++
		}
	}

	return rtn
}

func (o Offsets) ToMap() Map {
	if len(o) == 0 {
		return nil
	}

	rtn := make(Map, 0, len(o))

	rtn = append(rtn, Range{
		DestRangeStart: o[0].OffsetBy,
		SrcRangeStart:  0,
		Length:         o[0].InputTo + 1,
	})

	for i := 1; i < len(o); i++ {
		last := o[i-1]
		this := o[i]

		rangeStart := last.InputTo + 1

		rtn = append(rtn, Range{
			DestRangeStart: rangeStart + this.OffsetBy,
			SrcRangeStart:  rangeStart,
			Length:         this.InputTo - last.InputTo,
		})
	}

	return rtn
}

type Offset struct {
	// InputTo represents the inclusive upper range of this offset
	// The previous offsets InputTo +1 marks the start of this range
	//
	// If there is no previous offset, then the implicit start
	// is zero
	InputTo int

	// OffsetBy presents how much to offset the input
	// to create the output
	OffsetBy int
}

func (o Offset) String() string {
	if o.OffsetBy <= 0 {
		return fmt.Sprintf("Offset(%d, %d)", o.InputTo, o.OffsetBy)
	} else {
		return fmt.Sprintf("Offset(%d, +%d)", o.InputTo, o.OffsetBy)
	}
}

func (m Map) ToOffsets() (rtn Offsets, err error) {
	prev := Offset{InputTo: -1}

	for _, r := range m {
		// If this range starts with a gap
		// after the previous then insert a 0 offset
		// in the middle
		if r.SrcRangeStart > prev.InputTo+1 {
			prev = Offset{
				InputTo:  r.SrcRangeStart - 1,
				OffsetBy: 0,
			}
			rtn = append(rtn, prev)
		}

		if r.SrcRangeStart <= prev.InputTo {
			return nil, errors.Newf("range wasn't sorted")
		}

		prev = Offset{
			InputTo:  r.SrcRangeStart + r.Length - 1,
			OffsetBy: r.DestRangeStart - r.SrcRangeStart,
		}
		rtn = append(rtn, prev)
	}

	return rtn.Compress(), err
}

func MapsToMergedOffset(maps ...Map) (rtn Offsets, err error) {

	for _, m := range maps {
		newOffsets, err := m.ToOffsets()
		if err != nil {
			return nil, err
		}

		// On the first loop just set the rtn to the newOffsets
		if len(rtn) == 0 {
			rtn = newOffsets
			continue
		}

		newRtn := make(Offsets, 0, max(len(rtn)+len(newOffsets)))
		for i, currentOffset := range rtn {
			// Get the range start / end (inclsuive) for the current offset
			rangeStart := 0
			if i > 0 {
				rangeStart = rtn[i-1].InputTo + 1
			}
			rangeEnd := currentOffset.InputTo

			// Get the offset range start/end for the newOffsets
			newRangeStart := rangeStart + currentOffset.OffsetBy
			newRangeEnd := rangeEnd + currentOffset.OffsetBy

			// Find the index of the first newOffset contains the newRangeStart
			newOffsetIdx := newOffsets.FindIdx(newRangeStart)

			// If the newOffsetIdx is -1 then there's no newOffset which contains the newRangeStart
			// so we just add the currentOffset to the newRtn
			if newOffsetIdx == -1 {
				newRtn = append(newRtn, currentOffset)
				continue
			}

			for {
				if newOffsetIdx >= len(newOffsets) {
					newRtn = append(newRtn, Offset{
						InputTo:  rangeEnd,
						OffsetBy: currentOffset.OffsetBy,
					})
					break
				}

				newOffset := newOffsets[newOffsetIdx]

				// If the new offset ends after the current offset, then we can just add the current offset
				// with the newOffset's offsetBy added to it
				if newRangeEnd <= newOffset.InputTo {
					newRtn = append(newRtn, Offset{
						InputTo:  rangeEnd,
						OffsetBy: currentOffset.OffsetBy + newOffset.OffsetBy,
					})
					break
				}

				// Otherwise we need to work out the offset split point
				end := newOffset.InputTo - currentOffset.OffsetBy
				newRtn = append(newRtn, Offset{
					InputTo:  end,
					OffsetBy: currentOffset.OffsetBy + newOffset.OffsetBy,
				})
				newRangeStart = end + currentOffset.OffsetBy + 1
				newOffsetIdx++
			}
		}

		rtn = newRtn.Compress()
	}

	return rtn, nil
}
