//go:build !solution

package treeiter

func DoInOrder[e interface {
	Left() *e
	Right() *e
}](root *e, f func(t *e)) {

	if root == nil {
		return
	}

	DoInOrder((*root).Left(), f)
	f(root)
	DoInOrder((*root).Right(), f)
}
