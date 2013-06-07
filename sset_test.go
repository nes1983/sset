package sset

import "testing"

type intNode struct {
	val         int
	left, right *intNode
	color       Color
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

func TestHarmess(t *testing.T) {
	var set SortedSet
	_, found := Search(set, &intNode{val: 2})
	if found {
		t.Errorf("expecting to not find 2, but did.")
	} 

	Insert(&set, &intNode{val: 1})
	if l := Len(set); l != 1 {
		t.Errorf("expecting len 1, but got %d", l)
	}
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
	_, found = Search(set, &intNode{val: 2})
	if found {
		t.Errorf("expecting to not find 2, but did.")
	}
}

func BenchmarkInsert(b *testing.B) {
	var set SortedSet
	for i := 0; i < b.N; i++ {
		Insert(&set, &intNode{val: b.N - i})
	}
}
