package zone

import (
	"context"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/gdamore/tcell/v2"
	"github.com/mmta41/dnsimple-cli/commands/cmdutils"
	recordCmd "github.com/mmta41/dnsimple-cli/commands/zone/record"
	"github.com/mmta41/dnsimple-cli/pkg/iostreams"
	"github.com/mmta41/dnsimple-cli/pkg/prompt"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"log"
	"runtime/debug"
	"strconv"
)

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <zone>",
		Args:    cmdutils.MinimumArgs(1, "no zone name specified"),
		Short:   "View specific zone.",
		Long:    heredoc.Doc(`Show Records belong to specific zone`),
		PreRunE: f.PreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			return viewZone(f, args[0])
		},
	}

	return cmd
}

func viewZone(f *cmdutils.Factory, zoneName string) error {
	client, err := f.Client()
	if err != nil {
		return err
	}

	zoneRecords, err2 := client.Zones.ListRecords(context.Background(), client.AccountIdStr(), zoneName, &dnsimple.ZoneRecordListOptions{})
	if err2 != nil {
		return client.ConvertError(err2)
	}

	//Todo: Add GUI Output option
	if true {
		return cui(f, zoneRecords.Data, zoneName)
	}

	if f.IO.PromptEnabled() {
		return interactive(f, zoneRecords.Data, zoneName)
	}

	return text(f.IO, zoneRecords.Data)
}

func text(io *iostreams.IOStreams, zoneRecords []dnsimple.ZoneRecord) error {
	tHeader := fmt.Sprintf("ID\nName\tType\tTTL\tValue")
	fmt.Fprintf(io.StdOut, tHeader)
	for _, record := range zoneRecords {
		fmt.Fprintf(io.StdOut, "%v\t%v\t%v\t%v\t%v\n", record.ID, record.Name, record.Type, record.TTL, record.Content)
	}
	return nil
}

func interactive(f *cmdutils.Factory, zoneRecords []dnsimple.ZoneRecord, zoneName string) error {
	tHeader := fmt.Sprintf("ID\nName\tType\tTTL\tValue")
	var options = make([]string, len(zoneRecords))
	for i, record := range zoneRecords {
		options[i] = fmt.Sprintf("%v\t%v\t%v\t%v\t%v", record.ID, record.Name, record.Type, record.TTL, record.Content)
	}
	var selected int
	err := prompt.SelectOne(&selected, "Select and record:\n "+tHeader, options)
	if err != nil {
		return err
	}
	if selected >= 0 {
		return recordCmd.ViewRecord(f, zoneName, zoneRecords[selected].ID)
	}
	return nil
}

func cui(f *cmdutils.Factory, records []dnsimple.ZoneRecord, zoneName string) error {
	app := tview.NewApplication()
	recoverPanic(app)

	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)
	table.SetCell(0, 0, &tview.TableCell{
		Text:  "ID",
		Align: tview.AlignCenter,
		Color: tcell.ColorYellow,
	})
	table.SetCell(0, 1, &tview.TableCell{
		Text:  "Name",
		Align: tview.AlignCenter,
		Color: tcell.ColorYellow,
	})
	table.SetCell(0, 2, &tview.TableCell{
		Text:  "Type",
		Align: tview.AlignCenter,
		Color: tcell.ColorYellow,
	})
	table.SetCell(0, 3, &tview.TableCell{
		Text:  "TTL",
		Align: tview.AlignCenter,
		Color: tcell.ColorYellow,
	})
	table.SetCell(0, 4, &tview.TableCell{
		Text:  "Value",
		Align: tview.AlignCenter,
		Color: tcell.ColorYellow,
	})

	for rowId, record := range records {

		table.SetCell(rowId+1, 0, &tview.TableCell{
			Text:  strconv.FormatInt(record.ID, 10),
			Align: tview.AlignLeft,
			Color: tcell.ColorWhite,
		})
		table.SetCell(rowId+1, 1, &tview.TableCell{
			Text:  record.Name,
			Align: tview.AlignLeft,
			Color: tcell.ColorWhite,
		})
		table.SetCell(rowId+1, 2, &tview.TableCell{
			Text:  record.Type,
			Align: tview.AlignLeft,
			Color: tcell.ColorWhite,
		})
		table.SetCell(rowId+1, 3, &tview.TableCell{
			Text:  strconv.Itoa(record.TTL),
			Align: tview.AlignLeft,
			Color: tcell.ColorWhite,
		})
		table.SetCell(rowId+1, 4, &tview.TableCell{
			Text:  record.Content,
			Align: tview.AlignLeft,
			Color: tcell.ColorWhite,
		})

	}

	table.SetDoneFunc(func(key tcell.Key) {
	}).SetSelectedFunc(func(row, column int) {
		if row == 1 {
			return
		}
		app.Stop()
		recordCmd.ViewRecord(f, zoneName, records[row-1].ID)
	})
	return app.SetRoot(table, true).
		SetFocus(table).
		Run()
}

func recoverPanic(app *tview.Application) {
	if r := recover(); r != nil {
		app.Stop()
		log.Fatalf("%s\n%s\n", r, string(debug.Stack()))
	}
}
