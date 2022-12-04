package tree

//
// 学习链接: http://en.wikipedia.org/wiki/Rbtree
//
//  1) 所有节点要么红色要么黑色
//  2) root根节点是黑色。
//  3) 所有的叶子节点都是黑色，叶子节点全局统一定义nilLeaf
//  4) 红色节点的两个子节点是黑色
//  5) 从任一结点到其每个叶子的所有路径都包含相同数目的黑色结点
//

// Color 红黑树节点的颜色
type Color = int

const RED Color = 1
const BLACK Color = 2

type RbTreeKey interface {
	Less(than RbTreeKey) bool
	Equal(k RbTreeKey) bool
}

// RbTreeNode 红黑树节点
type RbTreeNode struct {
	Key    RbTreeKey
	Value  any
	color  Color
	left   *RbTreeNode
	right  *RbTreeNode
	parent *RbTreeNode
}

// NewRbTreeNode 生成新的节点
func (t *RbTree) newRbTreeNode(key RbTreeKey, val any, parent *RbTreeNode) *RbTreeNode {
	return &RbTreeNode{
		Key:    key,
		Value:  val,
		color:  RED,
		left:   t.null,
		right:  t.null,
		parent: parent,
	}
}

// RbTree 红黑树结构
type RbTree struct {
	null  *RbTreeNode
	root  *RbTreeNode
	count int
}

// NewRbTree 返回一个新的红黑树
func NewRbTree() *RbTree {
	node := &RbTreeNode{color: BLACK}
	return &RbTree{
		null:  node,
		root:  node,
		count: 0,
	}
}

// 红黑树左旋
func (t *RbTree) rbTreeLeftRotate(X *RbTreeNode) {
	// 左旋时X的右节点不能为空
	if X.right == t.null {
		return
	}

	// 左旋示意图
	//
	//          |                                  |
	//          X                                  Y
	//         / \            左旋                 / \
	//        A   Y      ------------->          X   C
	//           / \                            / \
	//          B   C                          A   B
	//
	// 旋转时不变色，之后统一处理变色问题

	Y := X.right
	X.right = Y.left
	if Y.left != t.null {
		Y.left.parent = X
	}

	Y.parent = X.parent
	if Y.parent == t.null {
		t.root = Y
	} else if Y.parent.left == X {
		Y.parent.left = Y
	} else if Y.parent.right == X {
		Y.parent.right = Y
	}

	Y.left = X
	X.parent = Y
}

// 红黑树右旋
func (t *RbTree) rbTreeRightRotate(X *RbTreeNode) {
	// 右旋时X的左节点不能为空
	if X.left == t.null {
		return
	}

	// 右旋示意图
	//
	//          |                                  |
	//          X                                  Y
	//         / \             右旋                / \
	//        Y   C      ------------->          A   X
	//       / \                                    / \
	//      A   B                                   B  C
	//
	// 旋转时不变色，之后统一处理变色问题

	Y := X.left
	X.left = Y.right
	if X.left != t.null {
		X.left.parent = X
	}

	Y.parent = X.parent
	if Y.parent == t.null {
		t.root = Y
	} else if Y.parent.left == X {
		Y.parent.left = Y
	} else if Y.parent.right == X {
		Y.parent.right = Y
	}

	Y.right = X
	X.parent = Y
}

// Update 修改节点
func (t *RbTree) Update(key RbTreeKey, val any) bool {
	curr := t.root
	for curr != t.null {
		if curr.Key.Less(key) {
			curr = curr.right
		} else if curr.Key.Equal(key) {
			curr.Value = val
			return true
		} else {
			curr = curr.left
		}
	}
	return false
}

// Insert 插入节点
func (t *RbTree) Insert(key RbTreeKey, val any) bool {
	prev := t.null
	curr := t.root

	for curr != t.null {
		prev = curr
		if curr.Key.Less(key) {
			curr = curr.right
		} else if curr.Key.Equal(key) {
			return false
		} else {
			curr = curr.left
		}
	}

	// 插入新节点初始化left，right字段，赋予红色，设置key，val，parent字段
	newNode := t.newRbTreeNode(key, val, prev)

	if prev == t.null {
		t.root = newNode
	} else if prev.Key.Less(newNode.Key) {
		prev.right = newNode
	} else {
		prev.left = newNode
	}

	t.count++
	// 平衡红黑树
	t.insertFixUp(newNode)
	return true
}

// 红黑树插入新节点后，可能会违背第四天和第五条约定
// insertFixUp 插入后重新调节红黑树
func (t *RbTree) insertFixUp(newNode *RbTreeNode) {
	// 父节点时黑色时，没有违反第四条和第五条约定，无需处理
	cur := newNode
	// cur节点会一直为红色，当父节点也是红色时，需要做出变色和旋转
	for cur.parent.color == RED {
		if cur.parent == cur.parent.parent.left { //父节点是祖父节点的左子树
			uncle := cur.parent.parent.right
			if uncle.color == RED { // 叔节点为红色
				// 场景1：父节点是祖父节点的左子树，且叔节点为红色，此时将父和叔节点变色为黑色，祖父节点变为红色，再去平衡祖父节点。
				cur.parent.color = BLACK
				uncle.color = BLACK
				cur.parent.parent.color = RED
				cur = cur.parent.parent
			} else { // 叔节点为黑色
				if cur == cur.parent.right {
					// 场景2：父节点是祖父节点的左子树，且叔节点为黑色，且当前节点是父节点的右子节点，此时围绕父节点进行左旋
					cur = cur.parent
					t.rbTreeLeftRotate(cur)
				}
				// 场景3：父节点是祖父节点的左子树，且叔节点为黑色，且当前节点是父节点的左子节点
				cur.parent.color = BLACK
				cur.parent.parent.color = RED
				t.rbTreeRightRotate(cur.parent.parent)
			}
		} else {
			uncle := cur.parent.parent.left
			if uncle.color == RED {
				cur.parent.color = BLACK
				uncle.color = BLACK
				cur.parent.parent.color = RED
				cur = cur.parent.parent
			} else {
				if cur == cur.parent.left {
					cur = cur.parent
					t.rbTreeRightRotate(cur)
				}
				cur.parent.color = BLACK
				cur.parent.parent.color = RED
				t.rbTreeLeftRotate(cur.parent.parent)
			}
		}
	}
	t.root.color = BLACK
	return
}

// Get 查询节点
func (t *RbTree) Get(key RbTreeKey) (any, bool) {
	curr := t.root
	for curr != t.null {
		if curr.Key.Less(key) {
			curr = curr.left
		} else if curr.Key.Equal(key) {
			return curr.Value, true
		} else {
			curr = curr.right
		}
	}
	return nil, false
}

func (t *RbTree) successor(x *RbTreeNode) *RbTreeNode {
	if x == t.null {
		return t.null
	}

	if x.right != t.null {
		return t.min(x.right)
	}

	y := x.parent
	for y != t.null && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

// Delete 删除节点
func (t *RbTree) Delete(key RbTreeKey) *RbTreeNode {
	curr := t.root
	for curr != t.null {
		if curr.Key.Less(key) {
			curr = curr.right
		} else if curr.Key.Equal(key) {
			break
		} else {
			curr = curr.left
		}
	}

	// 节点不存在，返回空
	if curr == t.null {
		return nil
	}

	// 节点存在时，进行删除操作，并返回节点
	ret := &RbTreeNode{
		Key:    curr.Key,
		Value:  curr.Value,
		color:  curr.color,
		left:   t.null,
		right:  t.null,
		parent: t.null,
	}

	var y *RbTreeNode
	var x *RbTreeNode
	if curr.left == t.null || curr.right == t.null {
		y = curr
	} else {
		y = t.successor(curr)
	}

	if y.left != t.null {
		x = y.left
	} else {
		x = y.right
	}

	x.parent = y.parent

	if y.parent == t.null {
		t.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != curr {
		curr.Key = y.Key
	}
	if y.color == BLACK {
		t.deleteFixUp(x)
	}

	t.count--
	return ret
}

// 删除节点后，进行平很
func (t *RbTree) deleteFixUp(x *RbTreeNode) {
	for x != t.root && x.color == BLACK {
		if x == x.parent.left {
			w := x.parent.right
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.rbTreeLeftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.right.color == BLACK {
					w.left.color = BLACK
					w.color = RED
					t.rbTreeRightRotate(w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.right.color = BLACK
				t.rbTreeLeftRotate(x.parent)
				// this is to exit while loop
				x = t.root
			}
		} else { // the code below is has left and right switched from above
			w := x.parent.left
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.rbTreeRightRotate(x.parent)
				w = x.parent.left
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					t.rbTreeLeftRotate(w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.left.color = BLACK
				t.rbTreeRightRotate(x.parent)
				x = t.root
			}
		}
	}
	x.color = BLACK
}

// Range DFS遍历该红黑树
func (t *RbTree) Range(f func(key RbTreeKey, val any)) {
	var dfs func(*RbTreeNode)
	dfs = func(cur *RbTreeNode) {
		if cur == t.null {
			return
		}
		dfs(cur.left)
		f(cur.Key, cur.Value)
		dfs(cur.right)
	}

	dfs(t.root)
}

// Min 获得最小key
func (t *RbTree) Min() *RbTreeNode {
	return t.min(t.root)
}

// 获得某个节点的最小子节点
func (t *RbTree) min(cur *RbTreeNode) *RbTreeNode {
	if cur == t.null {
		return t.null
	}

	cur = t.root
	for cur.left != t.null {
		cur = cur.left
	}
	return cur
}

// Max 获得最大key
func (t *RbTree) Max() *RbTreeNode {
	return t.max(t.root)
}

// 获得某个节点的最大子节点
func (t *RbTree) max(cur *RbTreeNode) *RbTreeNode {
	if cur == t.null {
		return t.null
	}

	cur = t.root
	for cur.right != t.null {
		cur = cur.right
	}
	return cur
}

func (t *RbTree) Size() int {
	return t.count
}
