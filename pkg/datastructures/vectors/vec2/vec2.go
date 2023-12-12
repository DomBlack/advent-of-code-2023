package vec2

type Vec2i [2]int

func (v Vec2i) Sub(v2 Vec2i) Vec2i {
	return Vec2i{v[0] - v2[0], v[1] - v2[1]}
}

func (v Vec2i) Length() int {
	x := v[0]
	if x < 0 {
		x = -x
	}
	y := v[1]
	if y < 0 {
		y = -y
	}
	return x + y
}
