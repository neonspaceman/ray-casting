package helpers

// ClampColor ensures a helpers value stays within 0-255 range
func ClampColor(value float32) uint8 {
	if value > 255 {
		return 255
	}

	if value < 0 {
		return 0
	}

	return uint8(value)
}
