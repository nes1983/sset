package sset

//import "log"

const (
	red   color = false // This must be the default.
	black color = true
)

type color bool

// Node must be implemented by nodes of the set.
type Node interface {
	// GetNodeInfo returns the node's set metadata.
	// It is usually implemented by embedding NodeInfo.
	GetNodeInfo() *NodeInfo

	// Cmp returns -1, if z < nd, 0 if z == nd, 1 if z > nd.
	Cmp(nd Node) int

	// SetValue updates the value to that of nd.
	SetValue(nd Node)
}

// NodeInfo holds a node's set metadata.
type NodeInfo struct {
	left, right Node
	color       color
}

func (info *NodeInfo) GetNodeInfo() *NodeInfo {
	return info
}

func nodeInfo(h Node) *NodeInfo {
	if h == nil {
		return nil
	}
	return h.GetNodeInfo()
}

// SortedSet holds a set of nodes.
type SortedSet struct {
	root Node
}

// Get returns a Node with the same value as nd,
// or nil if none was found.
func (set *SortedSet) Get(nd Node) Node {
	if set.root == nil {
		return nil
	}

	x := set.root
	for {
		cmp := x.Cmp(nd)
		if cmp == 0 {
			return x
		}
		info := x.GetNodeInfo()
		if cmp < 0 {
			x = info.left
		} else {
			x = info.right
		}
		if x == nil {
			return nil
		}
	}
}

// Len returns the number of elements in the set.
func (set *SortedSet) Len() int {
	return nodeLen(set.root)
}

func nodeLen(h Node) int {
	if h == nil {
		return 0
	}
	info := h.GetNodeInfo()
	return 1 + nodeLen(info.left) + nodeLen(info.right)
}

// Insert adds the given node to the set. If a node
// already exists with the same value, it will be
// changed to nd's value.
func (set *SortedSet) Insert(nd Node) {
	info := nd.GetNodeInfo()
	if info.left != nil || info.right != nil {
		panic("Nodes to be inserted should be nude, but weren't.")
	}
	set.root = insert(set.root, nd)
	set.root.GetNodeInfo().color = black
}

func insert(h Node, in Node) Node {
	if h == nil {
		return in
	}
	hinfo := h.GetNodeInfo()
	l, r := nodeInfo(hinfo.left), nodeInfo(hinfo.right)
	if l != nil && r != nil && l.color == red && r.color == red {
		colorFlip(hinfo, l, r)
	}
	if cmp := h.Cmp(in); cmp == 0 {
		h.SetValue(in)
	} else if cmp < 0 {
		hinfo.left = insert(hinfo.left, in)
		l = nodeInfo(hinfo.left)
	} else {
		hinfo.right = insert(hinfo.right, in)
		r = nodeInfo(hinfo.right)
	}

	if r != nil && r.color == red && !(l != nil && l.color == red) {
		h = rotateLeft(h, hinfo, r)
		hinfo = h.GetNodeInfo()
		l, r = nodeInfo(hinfo.left), nodeInfo(hinfo.right)
	}
	if l != nil && l.left != nil {
		ll := l.left.GetNodeInfo()
		if l.color == red && ll.color == red {
			h = rotateRight(h, hinfo, l)
		}
	}
	return h
}

func rotateLeft(h Node, hinfo, rinfo *NodeInfo) Node {
	x := hinfo.right
	hinfo.right = rinfo.left
	rinfo.left = h
	rinfo.color = hinfo.color
	hinfo.color = red
	return x
}

func rotateRight(h Node, hinfo, linfo *NodeInfo) Node {
	x := hinfo.left
	hinfo.left = linfo.right
	linfo.right = h
	linfo.color = hinfo.color
	hinfo.color = red
	return x
}

func colorFlip(hinfo, l, r *NodeInfo) {
	hinfo.color = !hinfo.color
	l.color = !l.color
	r.color = !r.color
}
