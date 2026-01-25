package example

import (
	"fmt"

	"github.com/l00pss/treego/rtree"
)

func main() {

	// (min: 2, max: 4 entry per node)
	tree := rtree.NewRTree(2, 4)

	item := &rtree.Item{
		Bounds: rtree.NewRectangle(0, 0, 10, 10),
		Data:   "Restaurant A",
	}
	tree.Insert(item)

	// search range
	results := tree.Search(rtree.NewRectangle(5, 5, 15, 15))
	fmt.Println(results)

	// search point
	point := rtree.Point{X: 7, Y: 8}
	items := tree.SearchPoint(point)
	fmt.Println(items)

	// find nearest 5 point
	nearest := tree.NearestNeighbor(point, 5)
	fmt.Println(nearest)
}
