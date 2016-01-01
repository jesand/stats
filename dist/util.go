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

// Find the minimum of an array of values
func Min(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	var min = x[0]
	for _, v := range x {
		if v < min {
			min = v
		}
	}
	return min
}

// Find the smallest value strictly greater than the provided value. If there
// is no such value, then val is returned.
func MinGt(x []float64, val float64) float64 {
	if len(x) == 0 {
		return 0
	}
	var min = val
	for _, v := range x {
		if val < v && (v < min || min == val) {
			min = v
		}
	}
	return min
}

// Find the maximum of an array of values
func Max(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	var max = x[0]
	for _, v := range x {
		if v > max {
			max = v
		}
	}
	return max
}

// Find the largest value strictly less than the provided value. If there
// is no such value, then val is returned.
func MaxLt(x []float64, val float64) float64 {
	if len(x) == 0 {
		return 0
	}
	var max = val
	for _, v := range x {
		if val > v && (v > max || max == val) {
			max = v
		}
	}
	return max
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
