package utils

import (
	"image/color"

	"github.com/cockroachdb/errors"
)

// ParseHexColor parses a hex colour code into a colour.RGBA
//
// The input is expected to be in the format: #RRGGBB or #RGB, #RRGGBBAA or #RGBA
func ParseHexColor(input string) (colour color.RGBA, err error) {
	if len(input) == 0 {
		return colour, errors.New("invalid hex code: empty string")
	}
	if input[0] != '#' {
		return colour, errors.Newf("invalid hex code; missing prefix: %q", input)
	}

	colour.A = 0xff

	decode := func(c byte) byte {
		switch {
		case c >= '0' && c <= '9':
			return c - '0'
		case c >= 'a' && c <= 'f':
			return c - 'a' + 0xa
		case c >= 'A' && c <= 'F':
			return c - 'A' + 0xa
		default:
			err = errors.Newf("invalid hex character: %c", c)
			return 0
		}
	}

	switch len(input) {
	case 9:
		colour.A = decode(input[7])<<4 + decode(input[8])
		fallthrough
	case 7:
		colour.R = decode(input[1])<<4 + decode(input[2])
		colour.G = decode(input[3])<<4 + decode(input[4])
		colour.B = decode(input[5])<<4 + decode(input[6])
	case 5:
		colour.A = decode(input[4]) * 0x11
		fallthrough
	case 4:
		colour.R = decode(input[1]) * 0x11
		colour.G = decode(input[2]) * 0x11
		colour.B = decode(input[3]) * 0x11
	default:
		return colour, errors.Newf("invalid hex code length: %q", input)
	}

	if err != nil {
		return colour, errors.Wrap(err, "invalid hex code")
	}
	return colour, nil
}
