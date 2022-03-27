package consolegame

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"goquiz/pkg/model"
)

const (
	menuPageNo gamePageNo = iota
	recordPageNo
	settingPageNo
)

type gamePageNo int
type setupPage func(pageNo gamePageNo, header string) tview.Primitive
type gamePage struct {
	pageNo gamePageNo
	setupPage
	header string
}

var currentPageNo gamePageNo = 0

var (
	app   = tview.NewApplication()
	pages = tview.NewPages()
)

var gamePages []gamePage

func getFrame(view tview.Primitive, header string) *tview.Frame {
	return tview.NewFrame(view).
		SetBorders(2, 2, 2, 2, 4, 4).
		AddText("[::b]"+header, true, tview.AlignLeft, tcell.ColorLightCyan).
		AddText("GoQuiz( early alpha ), credit [::b]https://javatpoint.com/[-:-:-] for mcq questions in this game", false, tview.AlignLeft, tcell.ColorLightCyan)
}

func getHeader(pageNo gamePageNo) string {
	for _, p := range gamePages {
		if p.pageNo == pageNo {
			return p.header
		}
	}
	return ""
}

func Start(q model.Quizzes) {
	gamePages = []gamePage{
		{
			pageNo:    menuPageNo,
			setupPage: menuPage(),
			header:    "MAIN MENU",
		},
		{
			pageNo:    recordPageNo,
			setupPage: recordPage(),
			header:    "GAME RECORD",
		},
	}

	for _, page := range gamePages {

		pageNo := page.pageNo
		header := page.header
		view := page.setupPage(pageNo, header)

		pages.AddPage(header, getFrame(view, header),
			true,
			currentPageNo == pageNo)
	}

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func switchGamePage(pageNo gamePageNo, header string) {
	currentPageNo = pageNo
	pages.SwitchToPage(header)
}
