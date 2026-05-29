package pkg

func DoubleToCents(value float64) int {
	return int(value * 100)
}

func CentsToDouble(cents int) float64 {
	return float64(cents) / 100
}
