package tree

import (
	"fmt"
	"testing"
)

type KEY struct {
	key int
}

func (v1 *KEY) Less(v2 RbTreeKey) bool {
	return v1.key < v2.(*KEY).key
}

func (v1 *KEY) Equal(v2 RbTreeKey) bool {
	return v1.key == v2.(*KEY).key
}

func TestNewRbTree2(t *testing.T) {
	tree := NewRbTree()
	tree.Insert(&KEY{key: 1}, 20)
	tree.Insert(&KEY{key: 10}, 1)
	tree.Insert(&KEY{key: 11}, 2)
	tree.Insert(&KEY{key: 17}, 3)
	tree.Insert(&KEY{key: -1}, 4)
	tree.Range(func(key RbTreeKey, val any) {
		fmt.Println(key.(*KEY).key, "  ", val)
	})
	fmt.Println(tree.Get(&KEY{key: 10}))
	fmt.Println(tree.Delete(&KEY{key: 1}).Key.(*KEY).key)
	tree.Range(func(key RbTreeKey, val any) {
		fmt.Println(key.(*KEY).key, "  ", val)
	})
}
