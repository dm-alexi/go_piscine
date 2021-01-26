package main

import (
	"container/heap"
)

// Present is a structure for evaluating presents
type Present struct {
	Value int
	Size  int
}

// PresentHeap is a max-heap of Presents (most valuable goes on top)
type PresentHeap []Present

// Len returns heap length
func (h PresentHeap) Len() int { return len(h) }

// Less provides PresentHeap ordering (most valuable goes on top)
func (h PresentHeap) Less(i, j int) bool {
	return h[i].Value > h[j].Value || (h[i].Value == h[j].Value && h[i].Size < h[j].Size)
}

// Swap swaps element in a PresentHeap
func (h PresentHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push present into a heap
func (h *PresentHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Present))
}

// Pop the most valuable present
func (h *PresentHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func getCoolestPresent(presents []Present) Present {
	h := &PresentHeap{}
	for _, v := range presents {
		heap.Push(h, v)
	}
	return heap.Pop(h).(Present)
}
