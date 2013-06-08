package sset

import "testing"
import "bytes"
import "fmt"

type intNode struct {
	NodeInfo
	val int
}

// Cmp returns -1, if z < nd, 0 if z == nd, 1 if z > nd.
func (z *intNode) Cmp(nd Node) int {
	return z.val - nd.(*intNode).val
}
func (z *intNode) SetValue(nd Node) {
	z.val = nd.(*intNode).val
}

func TestHarmless(t *testing.T) {
	var set SortedSet
	h := set.Get(&intNode{val: 2})
	if h != nil {
		t.Errorf("expecting to not find 2, but did.")
	}

	set.Insert(&intNode{val: int('a')})
	if l := set.Len(); l != 1 {
		t.Errorf("expecting len 1, but got %d", l)
	}
	set.Insert(&intNode{val: int('a')})
	if l := set.Len(); l != 1 {
		t.Errorf("expecting len 1, but got %d", l)
	}
	set.Insert(&intNode{val: int('b')})
	if l := set.Len(); l != 2 {
		t.Errorf("expecting len 2, but got %d", l)
	}
	actual := set.Get(&intNode{val: int('a')})
	if actual == nil {
		t.Errorf("expecting to find key, but was none.")
	} else if actual.(*intNode).val != int('a') {
		t.Errorf("Expecting to find key 1, but was %d", actual.(*intNode).val)
	}

	if actual := set.Get(&intNode{val: int('c')}); actual != nil {
		t.Errorf("Never added it, but still got %v.", actual)
	}
	set.Insert(&intNode{val: int('c')})
	set.Insert(&intNode{val: int('d')})
	set.Insert(&intNode{val: int('e')})

	if l := set.Len(); l != 5 {
		t.Errorf("Expecting len to be 5, but was %v", l)
	}

	checkTreeShape(t, describeTree(&set.root, false), "(a,(c,e)d)b;")
}

func checkTreeShape(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("set should have shape %v, but was %v", expected, actual)
	}
}

func TestReverseInsertHarmless(t *testing.T) {
	var set SortedSet
	set.Insert(&intNode{val: int('e')})
	set.Insert(&intNode{val: int('d')})
	set.Insert(&intNode{val: int('c')})
	checkTreeShape(t, describeTree(&set.root, false), "(c,e)d;")

	set.Insert(&intNode{val: int('b')})
	checkTreeShape(t, describeTree(&set.root, true), "((br,)cb,eb)db;")

	set.Insert(&intNode{val: int('a')})
	checkTreeShape(t, describeTree(&set.root, false), "((a,c)b,e)d;")
}

func BenchmarkInsert(t *testing.B) {
	var set SortedSet
	for i := 0; i < t.N; i++ {
		set.Insert(&intNode{val: t.N - i})
	}
}

func TestMakeAndDescribeTree(t *testing.T) {
	expected := "();"
	if actual := describeTree((*Node)(nil), false); actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
	for _, desc := range []string{
		"();",
		"((a,c)b,(e,g)f)d;",
	} {
		if actual := describeTree(makeTree(desc), false); actual != desc {
			t.Errorf("makeTree and describeTree should be symmetric, but weren't. "+
				"Input desc %v produced %v", desc, actual)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	//	b.StopTimer()
	var t SortedSet
	for i := 0; i < b.N; i++ {
		t.Insert(&intNode{val: b.N - i})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if i%1000 == 0 {
			//			log.Printf("%d", i)
		}
		t.Get(&intNode{val: i})
	}
}

// Build a tree from a simplified Newick format returning the root node.
// Single letter node names only, no error checking and all nodes are full or leaf.
func makeTree(desc string) *Node {
	if desc == "();" {
		return nil
	}
	buf := bytes.NewBufferString(desc)
	var build func() Node // Cannot be defined via ":=", because recursive.
	build = func() Node {
		if buf.Len() == 0 {
			return nil
		}

		var cn intNode
		for {
			b, _, _ := buf.ReadRune()
			switch b {
			case '(':
				cn.GetNodeInfo().left = build()
			case ',':
				cn.GetNodeInfo().right = build()
			case ')': // Ignore
			default:
				if b != ';' {
					cn.val = int(b)
				}
				return &cn
			}
		}

		panic("Unreachable")
	}

	n := build()
	if n.GetNodeInfo().left == nil && n.GetNodeInfo().right == nil {
		n = nil
	}

	return &n
}

// Return a Newick format description of a tree defined by a node.
func describeTree(n *Node, color bool) string {
	if n == nil {
		return "();"
	}

	var s bytes.Buffer

	var follow func(*Node)
	follow = func(n *Node) {
		info := (*n).GetNodeInfo()
		l, r := info.left, info.right
		if l != nil || r != nil {
			s.WriteString("(")
		}
		if l != nil {
			follow(&l)
		}
		if l != nil || r != nil {
			s.WriteString(",")
		}
		if r != nil {
			follow(&r)
		}
		if l != nil || r != nil {
			s.WriteString(")")
		}

		s.WriteString(fmt.Sprintf("%c", rune((*n).(*intNode).val)))
		if color {
			if (*n).GetNodeInfo().color == red {
				s.WriteRune('r')
			} else {
				s.WriteRune('b')
			}
		}
	}
	follow(n)

	s.WriteString(";")

	return s.String()
}

// Do all paths from root to leaf have same number of black edges?
func (t SortedSet) isBalanced() bool {
	if t.root == nil {
		return true
	}
	var nblack int // number of black links on path from root to min
	for x := t.root; x != nil; x = x.GetNodeInfo().left {
		l := x.GetNodeInfo().left
		if l != nil && l.GetNodeInfo().color == black {
			nblack++
		}
	}
	return isBalance(t.root.GetNodeInfo().left, nblack) &&
		isBalance(t.root.GetNodeInfo().right, nblack)
}

// Does every path from the root to a leaf have the given number
// of black links?
func isBalance(n Node, nblack int) bool {
	if n == nil && nblack == 0 {
		return true
	} else if n == nil && nblack != 0 {
		return false
	}
	if n.GetNodeInfo().color == black {
		nblack--
	}
	return isBalance(n.GetNodeInfo().left, nblack) &&
		isBalance(n.GetNodeInfo().right, nblack)
}

func TestRotateLeft(t *testing.T) {
	// Not sure how wise it is to test private methods ...
	spec := "((a,c)b,(e,g)f)d;"
	wanted := "(((a,c)b,e)d,g)f;"

	tree := *makeTree(spec)

	tree = rotateLeft(tree, tree.GetNodeInfo(), tree.GetNodeInfo().right.GetNodeInfo())
	if actual := describeTree(&tree, false); actual != wanted {
		t.Errorf("After rotation, wanted %v, but got %v.", actual, wanted)
	}
}

func TestInsertion(t *testing.T) {
	min := int('a')
	max := min + 1000
	var set SortedSet
	for i := min; i <= max; i++ {
		set.Insert(&intNode{val: i})

		if actual, expected := set.Len(), i-min+1; actual != expected {
			t.Errorf("Length should be %v, but was %v", expected, actual)
		}

		if !set.isBalanced() {
			t.Fatal("Tree should be balanced, but wasn't. Tree: ", describeTree(&set.root, true))
		}
		if !isBST(set.root, minimum(set.root), maximum(set.root)) {
			t.Fatal("Tree should be BST, but wasn't. Tree: ", describeTree(&set.root, false))
		}
		if !is234(set.root) {
			t.Fatal("Tree should be 23, but wasn't. Tree: ", describeTree(&set.root, true))
		}
	}
}

func TestGet(t *testing.T) {
	min := int('a')
	max := min + 1000
	var set SortedSet

	for i := min; i <= max; i++ {
		if i&1 == 0 {
			set.Insert(&intNode{val: i})
		}
	}

	for i := min; i <= max; i++ {
		desc := describeTree(&set.root, false)
		if i&1 == 0 {
			if actual := set.Get(&intNode{val: i}).(*intNode).val; actual != i {
				t.Errorf("Expected to find %v in set, but instead got %v. Tree: %v", i, actual, desc)
			}
		} else {
			if actual := set.Get(&intNode{val: i}); actual != nil {
				t.Errorf("Should not have found %v in set, but did. Tree: %v", i, desc)
			}
		}
	}
}

func is234(n Node) bool {
	if n == nil {
		return true
	}

	info := n.GetNodeInfo()
	if info.left != nil && info.right != nil &&
		info.left.GetNodeInfo().color == red && info.left.GetNodeInfo().color == black {
		return false
	}

	return is234(info.left) && is234(info.right)
}

// Are all the values in the BST rooted at x between min and max,
// and does the same property hold for both subtrees?
func isBST(n Node, min, max Node) bool {
	if n == nil {
		return true
	}

	if n.Cmp(min) < 0 || n.Cmp(max) > 0 {
		return false
	}

	return isBST(n.GetNodeInfo().left, min, n) && isBST(n.GetNodeInfo().right, n, max)
}

func minimum(n Node) Node {
	for l := n; l != nil; l = n.GetNodeInfo().left {
		n = l
	}
	return n
}

func maximum(n Node) Node {
	for l := n; l != nil; l = n.GetNodeInfo().right {
		n = l
	}
	return n
}
