package console

import (
	"github.com/rivo/tview"
)

func recordPage() setupPage {
	return func(pageNo gamePageNo, header string) tview.Primitive {

		list := tview.NewList().
			AddItem("Record A", " ", 'a', nil).
			AddItem("Record B", " ", 'b', nil).
			AddItem("Record C", " ", 'c', nil).
			AddItem("Back to Menu", "", 'd', func() {
				switchGamePage(pageNo, header)
			})
		return list
	}
}
