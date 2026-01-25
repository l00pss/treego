package rtree

import (
	"math"
)

// Point represents a point in 2D space
type Point struct {
	X, Y float64
}

// Rectangle represents a bounding box
type Rectangle struct {
	MinX, MinY, MaxX, MaxY float64
}

// Item represents an item to be stored in the R-tree
type Item struct {
	Bounds Rectangle
	Data   interface{}
}

// Node represents a node in the R-tree
type Node struct {
	isLeaf   bool
	bounds   Rectangle
	children []*Node
	items    []*Item
	parent   *Node
}

// RTree represents the R-tree structure
type RTree struct {
	root       *Node
	minEntries int
	maxEntries int
	size       int
}

// NewRTree creates a new R-tree with specified min/max entries per node
func NewRTree(minEntries, maxEntries int) *RTree {
	if minEntries < 1 || minEntries > maxEntries/2 {
		minEntries = maxEntries / 2
	}

	return &RTree{
		root:       &Node{isLeaf: true},
		minEntries: minEntries,
		maxEntries: maxEntries,
		size:       0,
	}
}

// NewRectangle creates a new rectangle
func NewRectangle(minX, minY, maxX, maxY float64) Rectangle {
	return Rectangle{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
	}
}

// NewPoint creates a point as a rectangle with zero area
func NewPoint(x, y float64) Rectangle {
	return Rectangle{x, y, x, y}
}

// Area calculates the area of a rectangle
func (r Rectangle) Area() float64 {
	return (r.MaxX - r.MinX) * (r.MaxY - r.MinY)
}

// Margin calculates the margin (perimeter) of a rectangle
func (r Rectangle) Margin() float64 {
	return (r.MaxX - r.MinX) + (r.MaxY - r.MinY)
}

// Intersects checks if two rectangles intersect
func (r Rectangle) Intersects(other Rectangle) bool {
	return r.MinX <= other.MaxX && r.MaxX >= other.MinX &&
		r.MinY <= other.MaxY && r.MaxY >= other.MinY
}

// Contains checks if this rectangle contains another
func (r Rectangle) Contains(other Rectangle) bool {
	return r.MinX <= other.MinX && r.MaxX >= other.MaxX &&
		r.MinY <= other.MinY && r.MaxY >= other.MaxY
}

// ContainsPoint checks if rectangle contains a point
func (r Rectangle) ContainsPoint(p Point) bool {
	return p.X >= r.MinX && p.X <= r.MaxX &&
		p.Y >= r.MinY && p.Y <= r.MaxY
}

// Expand expands this rectangle to include another
func (r *Rectangle) Expand(other Rectangle) {
	r.MinX = math.Min(r.MinX, other.MinX)
	r.MinY = math.Min(r.MinY, other.MinY)
	r.MaxX = math.Max(r.MaxX, other.MaxX)
	r.MaxY = math.Max(r.MaxY, other.MaxY)
}

// Union returns the smallest rectangle containing both rectangles
func (r Rectangle) Union(other Rectangle) Rectangle {
	return Rectangle{
		MinX: math.Min(r.MinX, other.MinX),
		MinY: math.Min(r.MinY, other.MinY),
		MaxX: math.Max(r.MaxX, other.MaxX),
		MaxY: math.Max(r.MaxY, other.MaxY),
	}
}

// EnlargementNeeded calculates area increase needed to include another rectangle
func (r Rectangle) EnlargementNeeded(other Rectangle) float64 {
	return r.Union(other).Area() - r.Area()
}

// Distance calculates minimum distance from rectangle to a point
func (r Rectangle) Distance(p Point) float64 {
	dx := math.Max(0, math.Max(r.MinX-p.X, p.X-r.MaxX))
	dy := math.Max(0, math.Max(r.MinY-p.Y, p.Y-r.MaxY))
	return math.Sqrt(dx*dx + dy*dy)
}

// Insert adds an item to the R-tree
func (t *RTree) Insert(item *Item) {
	t.size++
	leaf := t.chooseLeaf(t.root, item.Bounds)
	leaf.items = append(leaf.items, item)
	t.updateBounds(leaf)

	if len(leaf.items) > t.maxEntries {
		t.splitNode(leaf)
	}
}

// chooseLeaf finds the best leaf node to insert an item
func (t *RTree) chooseLeaf(node *Node, bounds Rectangle) *Node {
	if node.isLeaf {
		return node
	}

	var best *Node
	minEnlargement := math.MaxFloat64
	minArea := math.MaxFloat64

	for _, child := range node.children {
		enlargement := child.bounds.EnlargementNeeded(bounds)
		area := child.bounds.Area()

		if enlargement < minEnlargement ||
			(enlargement == minEnlargement && area < minArea) {
			minEnlargement = enlargement
			minArea = area
			best = child
		}
	}

	return t.chooseLeaf(best, bounds)
}

// updateBounds updates the bounding box of a node
func (t *RTree) updateBounds(node *Node) {
	if node.isLeaf {
		if len(node.items) == 0 {
			return
		}
		node.bounds = node.items[0].Bounds
		for i := 1; i < len(node.items); i++ {
			node.bounds.Expand(node.items[i].Bounds)
		}
	} else {
		if len(node.children) == 0 {
			return
		}
		node.bounds = node.children[0].bounds
		for i := 1; i < len(node.children); i++ {
			node.bounds.Expand(node.children[i].bounds)
		}
	}

	if node.parent != nil {
		t.updateBounds(node.parent)
	}
}

// splitNode splits an overflowing node using R*-tree splitting algorithm
func (t *RTree) splitNode(node *Node) {
	axis := t.chooseSplitAxis(node)
	index := t.chooseSplitIndex(node, axis)

	newNode := &Node{
		isLeaf: node.isLeaf,
		parent: node.parent,
	}

	if node.isLeaf {
		newNode.items = append([]*Item{}, node.items[index:]...)
		node.items = node.items[:index]
	} else {
		newNode.children = append([]*Node{}, node.children[index:]...)
		node.children = node.children[:index]
		for _, child := range newNode.children {
			child.parent = newNode
		}
	}

	t.updateBounds(node)
	t.updateBounds(newNode)

	if node.parent == nil {
		// Create new root
		t.root = &Node{
			isLeaf:   false,
			children: []*Node{node, newNode},
		}
		node.parent = t.root
		newNode.parent = t.root
		t.updateBounds(t.root)
	} else {
		node.parent.children = append(node.parent.children, newNode)
		if len(node.parent.children) > t.maxEntries {
			t.splitNode(node.parent)
		} else {
			t.updateBounds(node.parent)
		}
	}
}

// chooseSplitAxis determines the best axis to split on
func (t *RTree) chooseSplitAxis(node *Node) int {
	xMargin, yMargin := 0.0, 0.0

	if node.isLeaf {
		t.sortItemsByMinX(node.items)
		xMargin = t.calculateMarginSum(node, true)

		t.sortItemsByMinY(node.items)
		yMargin = t.calculateMarginSum(node, true)
	} else {
		t.sortNodesByMinX(node.children)
		xMargin = t.calculateMarginSum(node, false)

		t.sortNodesByMinY(node.children)
		yMargin = t.calculateMarginSum(node, false)
	}

	if xMargin < yMargin {
		return 0 // X axis
	}
	return 1 // Y axis
}

// calculateMarginSum calculates sum of margins for all distributions
func (t *RTree) calculateMarginSum(node *Node, isItems bool) float64 {
	sum := 0.0
	count := len(node.items)
	if !isItems {
		count = len(node.children)
	}

	for i := t.minEntries; i <= count-t.minEntries; i++ {
		r1, r2 := Rectangle{}, Rectangle{}

		if isItems {
			r1 = node.items[0].Bounds
			for j := 1; j < i; j++ {
				r1.Expand(node.items[j].Bounds)
			}
			r2 = node.items[i].Bounds
			for j := i + 1; j < count; j++ {
				r2.Expand(node.items[j].Bounds)
			}
		} else {
			r1 = node.children[0].bounds
			for j := 1; j < i; j++ {
				r1.Expand(node.children[j].bounds)
			}
			r2 = node.children[i].bounds
			for j := i + 1; j < count; j++ {
				r2.Expand(node.children[j].bounds)
			}
		}

		sum += r1.Margin() + r2.Margin()
	}

	return sum
}

// chooseSplitIndex determines the best index to split at
func (t *RTree) chooseSplitIndex(node *Node, axis int) int {
	count := len(node.items)
	if !node.isLeaf {
		count = len(node.children)
	}

	if axis == 0 {
		if node.isLeaf {
			t.sortItemsByMinX(node.items)
		} else {
			t.sortNodesByMinX(node.children)
		}
	} else {
		if node.isLeaf {
			t.sortItemsByMinY(node.items)
		} else {
			t.sortNodesByMinY(node.children)
		}
	}

	minOverlap := math.MaxFloat64
	minArea := math.MaxFloat64
	bestIndex := t.minEntries

	for i := t.minEntries; i <= count-t.minEntries; i++ {
		r1, r2 := Rectangle{}, Rectangle{}

		if node.isLeaf {
			r1 = node.items[0].Bounds
			for j := 1; j < i; j++ {
				r1.Expand(node.items[j].Bounds)
			}
			r2 = node.items[i].Bounds
			for j := i + 1; j < count; j++ {
				r2.Expand(node.items[j].Bounds)
			}
		} else {
			r1 = node.children[0].bounds
			for j := 1; j < i; j++ {
				r1.Expand(node.children[j].bounds)
			}
			r2 = node.children[i].bounds
			for j := i + 1; j < count; j++ {
				r2.Expand(node.children[j].bounds)
			}
		}

		overlap := 0.0
		if r1.Intersects(r2) {
			ix := math.Min(r1.MaxX, r2.MaxX) - math.Max(r1.MinX, r2.MinX)
			iy := math.Min(r1.MaxY, r2.MaxY) - math.Max(r1.MinY, r2.MinY)
			overlap = ix * iy
		}

		area := r1.Area() + r2.Area()

		if overlap < minOverlap || (overlap == minOverlap && area < minArea) {
			minOverlap = overlap
			minArea = area
			bestIndex = i
		}
	}

	return bestIndex
}

// Sorting functions
func (t *RTree) sortItemsByMinX(items []*Item) {
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Bounds.MinX > items[j].Bounds.MinX {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func (t *RTree) sortItemsByMinY(items []*Item) {
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Bounds.MinY > items[j].Bounds.MinY {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func (t *RTree) sortNodesByMinX(nodes []*Node) {
	for i := 0; i < len(nodes)-1; i++ {
		for j := i + 1; j < len(nodes); j++ {
			if nodes[i].bounds.MinX > nodes[j].bounds.MinX {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
		}
	}
}

func (t *RTree) sortNodesByMinY(nodes []*Node) {
	for i := 0; i < len(nodes)-1; i++ {
		for j := i + 1; j < len(nodes); j++ {
			if nodes[i].bounds.MinY > nodes[j].bounds.MinY {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
		}
	}
}

// Search finds all items that intersect with the given rectangle
func (t *RTree) Search(bounds Rectangle) []*Item {
	result := []*Item{}
	t.searchNode(t.root, bounds, &result)
	return result
}

func (t *RTree) searchNode(node *Node, bounds Rectangle, result *[]*Item) {
	if !node.bounds.Intersects(bounds) {
		return
	}

	if node.isLeaf {
		for _, item := range node.items {
			if item.Bounds.Intersects(bounds) {
				*result = append(*result, item)
			}
		}
	} else {
		for _, child := range node.children {
			t.searchNode(child, bounds, result)
		}
	}
}

// SearchPoint finds all items that contain the given point
func (t *RTree) SearchPoint(p Point) []*Item {
	result := []*Item{}
	t.searchPointNode(t.root, p, &result)
	return result
}

func (t *RTree) searchPointNode(node *Node, p Point, result *[]*Item) {
	if !node.bounds.ContainsPoint(p) {
		return
	}

	if node.isLeaf {
		for _, item := range node.items {
			if item.Bounds.ContainsPoint(p) {
				*result = append(*result, item)
			}
		}
	} else {
		for _, child := range node.children {
			t.searchPointNode(child, p, result)
		}
	}
}

// NearestNeighbor finds the k nearest items to a point
func (t *RTree) NearestNeighbor(p Point, k int) []*Item {
	type queueItem struct {
		node     *Node
		item     *Item
		distance float64
	}

	queue := []queueItem{{node: t.root, distance: t.root.bounds.Distance(p)}}
	result := []*Item{}

	for len(queue) > 0 && len(result) < k {
		// Find minimum distance item in queue
		minIdx := 0
		for i := 1; i < len(queue); i++ {
			if queue[i].distance < queue[minIdx].distance {
				minIdx = i
			}
		}

		current := queue[minIdx]
		queue = append(queue[:minIdx], queue[minIdx+1:]...)

		if current.item != nil {
			result = append(result, current.item)
			continue
		}

		if current.node.isLeaf {
			for _, item := range current.node.items {
				dist := item.Bounds.Distance(p)
				queue = append(queue, queueItem{item: item, distance: dist})
			}
		} else {
			for _, child := range current.node.children {
				dist := child.bounds.Distance(p)
				queue = append(queue, queueItem{node: child, distance: dist})
			}
		}
	}

	return result
}

// Size returns the number of items in the tree
func (t *RTree) Size() int {
	return t.size
}

// Height returns the height of the tree
func (t *RTree) Height() int {
	return t.getHeight(t.root)
}

func (t *RTree) getHeight(node *Node) int {
	if node.isLeaf {
		return 1
	}
	if len(node.children) == 0 {
		return 1
	}
	return 1 + t.getHeight(node.children[0])
}
