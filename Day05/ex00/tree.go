package main

import (
	"fmt"
	"strings"
)

// TreeNode is a boolean binary tree node
type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func newNode(toy bool) *TreeNode {
	return &TreeNode{toy, nil, nil}
}

func countToys(root *TreeNode) int {
	if root == nil {
		return 0
	}
	s := 0
	if root.HasToy {
		s = 1
	}
	return s + countToys(root.Left) + countToys(root.Right)
}

func areToysBalanced(root *TreeNode) bool {
	return countToys(root.Left) == countToys(root.Right)
}

func getHeight(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return max(getHeight(root.Left), getHeight(root.Right)) + 1
}

func fillLevels(root *TreeNode, levels [][]byte, level int) {
	if root == nil {
		tmp := []byte{' '}
		for i := level; i < len(levels); i++ {
			levels[i] = append(levels[i], tmp...)
			tmp = append(tmp, tmp...)
		}
		return
	}
	t := byte('0')
	if root.HasToy {
		t = '1'
	}
	levels[level] = append(levels[level], t)
	fillLevels(root.Left, levels, level+1)
	fillLevels(root.Right, levels, level+1)
}

func printTree(root *TreeNode) {
	height := getHeight(root)
	if height == 0 {
		return
	}
	levels := make([][]byte, height)
	fillLevels(root, levels, 0)
	for i, str := range levels {
		var sb strings.Builder
		sb.WriteString(strings.Repeat(" ", 1<<(height-i-1)-1))
		for _, ch := range str {
			sb.WriteByte(ch)
			sb.WriteString(strings.Repeat(" ", 1<<(height-i)-1))
		}
		fmt.Println(sb.String())
	}
}
