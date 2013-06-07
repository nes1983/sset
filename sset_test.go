package sset

import "testing"
import "bytes"
import "fmt"

type intNode struct {
	NodeInfo
	val int
}

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

	set.Insert(&intNode{val: 1})
	if l := set.Len(); l != 1 {
		t.Errorf("expecting len 1, but got %d", l)
	}
	set.Insert(&intNode{val: 1})
	if l := set.Len(); l != 1 {
		t.Errorf("expecting len 1, but got %d", l)
	}
	set.Insert(&intNode{val: 2})
	if l := set.Len(); l != 2 {
		t.Errorf("expecting len 2, but got %d", l)
	}
	actual := set.Get(&intNode{val: 1})
	if actual == nil {
		t.Errorf("expecting to find key, but was none.")
	} else if actual.(*intNode).val != 1 {
		t.Errorf("Expecting to find key 1, but was %d", actual.(*intNode).val)
	}

	if actual := set.Get(&intNode{val: 3}); actual != nil {
		t.Errorf("expecting to not find 2, but found %v.", actual)
	}
	set.Insert(&intNode{val: 3})
	set.Insert(&intNode{val: 4})
	set.Insert(&intNode{val: 5})

	if l := set.Len(); l != 5 {
		t.Errorf("Expecting len to be 5, but was %v", l)
	}
}

func BenchmarkInsert(b *testing.B) {
	var t SortedSet
	for i := 0; i < b.N; i++ {
		t.Insert(&intNode{val: b.N - i})
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
	min, max := 1, 1000
	var set SortedSet
	for i := min; i <= max; i++ {
		set.Insert(&intNode{val: i})
		if actual := set.Len(); actual != i {
			t.Errorf("Length should be %v, but was %v", i, actual)
		}
		if !set.isBalanced() {
			t.Errorf("Length: %v", set.Len())
			t.Fatalf("Tree should be balanced, but wasn't. Tree: %v", describeTree(&set.root, true))
		}
	}
	// TODO:
	//		failed = failed || !c.Check(t.isBST(), check.Equals, true)
	//		failed = failed || !c.Check(t.is23_234(), check.Equals, true)
}
