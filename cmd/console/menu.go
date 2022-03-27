package console

import (
	"fmt"
	"github.com/rivo/tview"
)

func menuPage() setupPage {
	return func(pageNo gamePageNo, header string) tview.Primitive {
		list := tview.NewList().
			AddItem("Play Game", " ", 'a', nil).
			AddItem("Game Records", " ", 'b', func() {
				switchGamePage(recordPageNo, getHeader(recordPageNo))
			}).
			AddItem("Setting", " ", 'c', nil).
			AddItem("Quit", "Press to exit", 'q', func() {
				app.Stop()
				fmt.Println(currentPageNo)
			})
		return list
	}

}
