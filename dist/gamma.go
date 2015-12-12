package dist

import (
	"math"
	"math/rand"
)

// Return a random value drawn from a Gamma distribution with mean
// alpha*beta+lamba and variance alpha*beta^2.
// Based on nextGamma() in Factorie: https://github.com/factorie/factorie
func randGamma(alpha, beta, lambda float64) float64 {
	var gamma float64
	if alpha <= 0 || beta <= 0 {
		panic(Errorf("Invalid Gamma distribution parameters: alpha=%f, beta=%f",
			alpha, beta))
	}
	if alpha < 1 {
		var (
			p = 0.0
			b = 1 + alpha*math.Exp(-1)
		)
		for {
			p = b * rand.Float64()
			if p > 1 {
				gamma = -math.Log((b - p) / alpha)
				if rand.Float64() <= math.Pow(gamma, alpha-1) {
					break
				}
			} else {
				gamma = math.Pow(p, 1/alpha)
				if rand.Float64() <= math.Exp(-gamma) {
					break
				}
			}
		}
	} else if alpha == 1 {
		gamma = -math.Log(rand.Float64())
	} else {
		var y = -math.Log(rand.Float64())
		for rand.Float64() > math.Pow(y*math.Exp(1-y), alpha-1) {
			y = -math.Log(rand.Float64())
		}
		gamma = alpha * y
	}
	return beta*gamma + lambda
}
