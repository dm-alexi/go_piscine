package distribution

import "math/rand"

// NormDist contains normal distribution parameters
type NormDist struct {
	Mean   float64
	Stddev float64
}

// GetDistribution returns NormDist with mean in [-10, 10], stddev in [0.3, 1.5]
func GetDistribution() *NormDist {
	return &NormDist{Mean: rand.Float64()*20 - 10, Stddev: rand.Float64()*1.2 + 0.3}
}

// GetEntry return a sample distributed according to dist parameter
func GetEntry(dist *NormDist) float64 {
	return rand.NormFloat64()*dist.Stddev + dist.Mean
}
