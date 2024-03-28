package treeiter_test

import (
	"fmt"

	"gitlab.com/manytask/itmo-go/public/treeiter"
)

func ExampleDoInOrder() {
	tree := &ValuesNode[string]{
		value: "root",
		left: &ValuesNode[string]{
			value: "left",
		},
		right: &ValuesNode[string]{
			value: "right",
		},
	}

	treeiter.DoInOrder(tree, func(t *ValuesNode[string]) {
		fmt.Println(t.value)
	})

	// Output:
	// left
	// root
	// right
}
