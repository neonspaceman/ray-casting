package helpers

type Color uint32

func (c Color) RGBA() (r, g, b, a uint32) {
	clr := uint32(c)
	r = clr >> 24
	r |= r << 8
	g = (clr >> 16) & 0xFF
	g |= g << 8
	b = (clr >> 8) & 0xFF
	b |= b << 8
	a = clr & 0xFF
	a |= a << 8
	return
}
