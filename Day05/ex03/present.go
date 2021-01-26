package main

// Present is a structure for evaluating presents
type Present struct {
	Value int
	Size  int
}

func grabPresents(presents []Present, capacity int) []Present {
	n := len(presents)
	m := make([][]int, n+1)
	for i := range m {
		m[i] = make([]int, capacity+1)
	}
	p := make([][][]Present, n+1)
	for i := range p {
		p[i] = make([][]Present, capacity+1)
		for j := range p[i] {
			p[i][j] = make([]Present, 0)
		}
	}
	for i := 1; i <= n; i++ {
		for j := 0; j <= capacity; j++ {
			if presents[i-1].Size > j || m[i-1][j] >= m[i-1][j-presents[i-1].Size]+presents[i-1].Value {
				m[i][j] = m[i-1][j]
				p[i][j] = p[i-1][j]
			} else {
				m[i][j] = m[i-1][j-presents[i-1].Size] + presents[i-1].Value
				p[i][j] = append(p[i-1][j-presents[i-1].Size], presents[i-1])
			}
		}
	}
	return p[n][capacity]
}
