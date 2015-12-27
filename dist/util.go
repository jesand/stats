package dist

// Compute the sum of an array of values
func Sum(x []float64) float64 {
	var total float64
	for _, v := range x {
		total += v
	}
	return total
}

// Compute the sample mean of an array of values
func Mean(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	return Sum(x) / float64(len(x))
}

// Compute the sample variance of an array of values
func Variance(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	var (
		sum  = Sum(x)
		mean = sum / float64(len(x))

		total float64
	)
	for _, v := range x {
		diff := mean - v
		total += diff * diff
	}
	return total / (float64(len(x)) - 1)
}
