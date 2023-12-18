package day15

import (
	"testing"
)

func Test_Day15(t *testing.T) {
	input := "rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7"

	Day15.Test(t, input, 1320, input, 145)
}
