package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sum-project/pass-locker2/db/repository"
	"golang.design/x/clipboard"
	"strconv"
)

func listPage(a *App) fyne.CanvasObject {
	passwords, err := a.passwordRepository.GetAllPasswords(repository.GetAllPasswordsParams{
		UserId: a.user.ID,
	})
	if err != nil {
		a.errorLog.Fatalln(err)
	}

	slice := passwordsSlice(passwords)

	table := widget.NewTable(
		func() (int, int) {
			return len(slice), len(slice[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewLabel(""))
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			if id.Col == (len(slice[0])-2) && id.Row != 0 {
				copyButton := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
					clipboard.Write(clipboard.FmtText, []byte(slice[id.Row][3].(string)))
				})
				copyButton.Importance = widget.HighImportance
				object.(*fyne.Container).Objects = []fyne.CanvasObject{
					copyButton,
				}
			} else if id.Col == (len(slice[0])-1) && id.Row != 0 {
				deleteButton := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Delete", "", func(deleted bool) {
						if deleted {
							id, _ := strconv.Atoi(slice[id.Row][0].(string))
							a.passwordRepository.DeletePassword(repository.DeletePasswordParams{
								ID: uint64(id),
							})
						}
					}, a.Win)
				})
				deleteButton.Importance = widget.HighImportance

				object.(*fyne.Container).Objects = []fyne.CanvasObject{
					deleteButton,
				}
			} else {
				object.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(slice[id.Row][id.Col].(string)),
				}
			}
		},
	)

	colWidths := []float32{50, 200, 200, 200, 110, 110}
	for i := 0; i < len(colWidths); i++ {
		table.SetColumnWidth(i, colWidths[i])
	}

	return container.NewAdaptiveGrid(1, table)
}

func passwordsSlice(passwords []repository.Password) [][]interface{} {
	var slice [][]interface{}

	slice = append(slice, []interface{}{"ID", "Item", "Login", "Password", "Copy", "Delete"})

	for _, x := range passwords {
		var row []interface{}

		row = append(row, strconv.FormatUint(uint64(x.ID), 10))
		row = append(row, x.Item)
		row = append(row, x.Login)
		row = append(row, x.Password)
		row = append(row, widget.NewButton("Copy", func() {}))
		row = append(row, widget.NewButton("Delete", func() {}))

		slice = append(slice, row)
	}

	return slice
}
