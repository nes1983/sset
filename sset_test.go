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

// Return a Newick format description of a tree defined by a node.
func describeTree(n *Node) string {
	if n == nil {
		return "()"
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
		col := func (c color) string { if c == red { return "r" }; return "b"}
		s.WriteString(fmt.Sprintf("%d %v", (*n).(*intNode).val, col((*n).GetNodeInfo().color)))
	}
	
	follow(n)
	
	s.WriteString(";")

	return s.String()
}
