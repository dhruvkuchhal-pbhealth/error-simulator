package handlers

import (
	"net/http"
)

// TreeNode represents a node in a binary tree.
type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
	Value int
}

// TreeOps provides tree operations. The bug: CalculateDepth recurses without
// a base case for nil, and we pass a node that points to itself so recursion never ends.
type TreeOps struct{}

// CalculateDepth returns the depth of the tree. When node is nil it should return 0.
// BUG: Missing base case — when node is nil we still recurse (e.g. node.Left) instead
// of returning 0, and with a circular reference the stack overflows.
func (t *TreeOps) CalculateDepth(node *TreeNode) int {
	// Should be: if node == nil { return 0 }
	return 1 + max(t.CalculateDepth(node.Left), t.CalculateDepth(node.Right))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// StackOverflow handles GET /error/stack-overflow.
// It builds a TreeNode that points to itself and calls CalculateDepth, causing
// infinite recursion and "fatal error: stack overflow". This cannot be recovered
// by the recovery middleware.
func StackOverflow(ops *TreeOps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		node := &TreeNode{Value: 1}
		node.Left = node // circular reference
		node.Right = node
		_ = ops.CalculateDepth(node)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
