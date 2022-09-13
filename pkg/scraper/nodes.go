package scraper

import "github.com/PuerkitoBio/goquery"

type nodeDirection int8

const (
	D_Next nodeDirection = iota // find in the next nodes from the current node
	D_Prev                      // find in the previous nodes from the current node
)

// find the sibling node which className is equal to the class provided in the argument.
// the search fould be forward or backward based on the direction
func findSibling(node *goquery.Selection, class string, depth int, dir nodeDirection) *goquery.Selection {

	for i := 1; i < depth; i++ {
		switch dir {
		case D_Next:
			node = node.Next()
		case D_Prev:
			node = node.Prev()
		default:
			panic(" Unknown Node Direction")
		}

		nodeclass, exist := node.Attr("class")
		if exist && nodeclass == class {
			return node
		}
	}
	return nil
}
