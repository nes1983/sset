package sset

// Copyright (C) 2013 Niko Schwarz 
// Free Software. There should be a COPYING file around.
// If not, see <http://www.gnu.org/licenses/>. 

const (
	RED   Color = true
	BLACK Color = false
)

type Color bool

type Node interface {
	// -1, if z < nd, 0 if z == nd, 1 if z > nd.
	Cmp(nd Node) int
	Left() (nd Node, ok bool)
	Right() (nd Node, ok bool)
	SetLeft(Node)
	SetRight(Node)
	Color() Color
	SetColor(Color)
	// Update the value to that of nd. Don't change Left or Right.
	SetValue(nd Node)
}

type SortedSet struct {
	root *Node
}

func Search(set SortedSet, nd Node) (Node, bool) {
	if set.root == nil {
		return nil, false
	}

	x := *set.root
	for {
		var ok bool
		cmp := x.Cmp(nd)
		switch {
		case cmp == 0:
			return x, true
		case cmp < 0:
			x, ok = x.Left()
		case cmp > 0:
			x, ok = x.Right()
		}
		if !ok {
			return nil, false
		}
	}
}

func Len(set SortedSet) int {
	if set.root == nil {
		return 0
	}
	return nodeLen(*set.root)
}

func nodeLen(h Node) (ret int) {
	ret = 1
	if l, ok := h.Left(); ok {
		ret += nodeLen(l)
	}
	if r, ok := h.Right(); ok {
		ret += nodeLen(r)
	}
	return
}

func Insert(set *SortedSet, nd Node) {
	_, leftOk := nd.Left()
	_, rightOk := nd.Right()
	if leftOk || rightOk {
		panic("Nodes to be inserted should be nude, but weren't.")
	}
	if set.root == nil {
		set.root = &nd
	} else {
		i := insert(*set.root, nd)
		set.root = &i // Must be on different lines: can't take pointer of ret val.
	}
	(*set.root).SetColor(BLACK)
}

func insert(h Node, in Node) Node {
	l, okL := h.Left()
	r, okR := h.Right()

	if okL && okR && l.Color() == RED && r.Color() == RED {
		colorFlip(h)
	}
	if cmp := in.Cmp(h); cmp == 0 {
		h.SetValue(in)
	} else if okL && cmp < 0 {
		h.SetLeft(insert(l, in))
	} else if okR {
		h.SetRight(insert(r, in))
	}

	if okL && okR && r.Color() == RED && !l.Color() == RED {
		h = rotateLeft(h)
	}
	if okL {
		if ll, okLL := l.Left(); okLL && l.Color() == RED && ll.Color() == RED {
			h = rotateRight(h)
		}
	}
	return h
}

func rotateLeft(h Node) (x Node) {
	var okX bool
	x, okX = h.Right()
	l, okL := x.Left()
	if !okX || !okL {
		panic("For rotation, children must be there")
	}

	h.SetRight(l)
	x.SetLeft(h)
	x.SetColor(h.Color())
	h.SetColor(RED)
	return
}

func rotateRight(h Node) (x Node) {
	var okX bool
	x, okX = h.Left()
	r, okR := x.Right()
	if !okX || !okR {
		panic("For rotation, children must be there")
	}

	h.SetLeft(r)
	x.SetRight(h)
	x.SetColor(h.Color())
	h.SetColor(RED)
	return
}

func colorFlip(h Node) {
	r, okR := h.Right()
	l, okL := h.Left()
	if !okL || !okR {
		panic("For rotation, children must be there")
	}
	h.SetColor(!h.Color())
	l.SetColor(!l.Color())
	r.SetColor(!r.Color())
}
