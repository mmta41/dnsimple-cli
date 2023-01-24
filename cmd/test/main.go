package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"runtime/debug"
)

func main() {
	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)
	rows := [][]string{
		{"Name", "TYPE", "TTL", "Value"},
		{"Ali", "10", "12", "232"},
		{"Hasan", "CX", "2444", "Gasnareroo"},
	}

	for rowId, row := range rows {
		for colId, col := range row {

			color := tcell.ColorWhite
			align := tview.AlignLeft

			if rowId == 0 {
				color = tcell.ColorYellow
				align = tview.AlignCenter
			}
			table.SetCell(rowId, colId, &tview.TableCell{
				Text:  col,
				Align: align,
				Color: color,
			})
		}
	}
	a := tview.NewApplication().SetRoot(table, true)
	a.SetFocus(table)
	a.Run()
}
func recoverPanic(app *tview.Application) {
	if r := recover(); r != nil {
		app.Stop()
		log.Fatalf("%s\n%s\n", r, string(debug.Stack()))
	}
}
