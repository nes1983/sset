sset, SortedSet, implemented as a left-leaning red-black tree (LLRB)
====

Based Left-leaning Red-Black Trees, by Robert Sedgewick, doi://10.1.1.139.282.

For an example of how to use it, see the unit test. In a nutshell:

	type intNode struct {
		val		 int
		left, right *intNode
		color	   Color
	}

	// -1, if z < nd, 0 if z == nd, 1 if z > nd.
	func (z intNode) Cmp(nd Node) int {
		v := nd.(*intNode).val
		if z.val < v {
			return -1
		} else if z.val == v {
			return 0
		} else {
			return 1
		}
	}

	func (z intNode) Left() (nd Node, ok bool)  { return z.left, z.left != nil }
	func (z intNode) Right() (nd Node, ok bool) { return z.right, z.right != nil }
	func (z *intNode) SetLeft(nd Node)          { (*z).left = nd.(*intNode) }
	func (z *intNode) SetRight(nd Node)         { (*z).right = nd.(*intNode) }
	func (z intNode) Color() Color              { return z.color }
	func (z *intNode) SetColor(cl Color)        { (*z).color = cl }
	func (z *intNode) SetValue(nd Node)         { z.val = nd.(*intNode).val }

	func HarmlessTest(t *testing.T) {
		var set SortedSet

		Insert(&set, &intNode{val: 1})
		if l := Len(set); l != 1 {
			t.Errorf("expecting len 1, but got %d", l)
		}
	
		actual, found := Search(set, &intNode{val: 1})
		if !found {
			t.Errorf("expecting to find key, but was none.")
		} else if actual.(*intNode).val != 1 {
			t.Errorf("Expecting to find key 1, but was %d", actual.(*intNode).val)
		}
	}
