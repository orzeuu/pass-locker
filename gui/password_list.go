package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sum-project/pass-locker2/db/repository"
	"github.com/sum-project/pass-locker2/password"
	"golang.design/x/clipboard"
	"strconv"
)

func (a *App) passwordListPage() fyne.CanvasObject {
	passwords, err := a.passwordRepository.GetAllPasswords(repository.GetAllPasswordsParams{
		UserId: a.user.ID,
	})
	if err != nil {
		a.errorLog.Fatalln(err)
	}

	a.passwords = a.passwordsSlice(passwords)

	a.passwordTable = widget.NewTable(
		func() (int, int) {
			return len(a.passwords), len(a.passwords[0])
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewLabel(""))
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			if id.Col == (len(a.passwords[0])-2) && id.Row != 0 {
				copyButton := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
					clipboard.Write(clipboard.FmtText, []byte(a.passwords[id.Row][3].(string)))
				})
				copyButton.Importance = widget.HighImportance
				object.(*fyne.Container).Objects = []fyne.CanvasObject{
					copyButton,
				}
			} else if id.Col == (len(a.passwords[0])-1) && id.Row != 0 {
				deleteButton := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Delete", "", func(deleted bool) {
						if deleted {
							id, _ := strconv.Atoi(a.passwords[id.Row][0].(string))
							a.passwordRepository.DeletePassword(repository.DeletePasswordParams{
								ID: uint64(id),
							})

							passwords, err = a.passwordRepository.GetAllPasswords(repository.GetAllPasswordsParams{
								UserId: a.user.ID,
							})
							if err != nil {
								a.errorLog.Fatalln(err)
							}

							a.passwords = a.passwordsSlice(passwords)
							a.passwordTable.Refresh()
						}
					}, a.Win)
				})
				deleteButton.Importance = widget.HighImportance

				object.(*fyne.Container).Objects = []fyne.CanvasObject{
					deleteButton,
				}
			} else {
				object.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(a.passwords[id.Row][id.Col].(string)),
				}
			}
		},
	)

	colWidths := []float32{50, 200, 200, 200, 110, 110}
	for i := 0; i < len(colWidths); i++ {
		a.passwordTable.SetColumnWidth(i, colWidths[i])
	}

	return container.NewAdaptiveGrid(1, a.passwordTable)
}

func (a *App) passwordsSlice(passwords []repository.Password) [][]interface{} {
	var slice [][]interface{}

	slice = append(slice, []interface{}{"ID", "Item", "Login", "Password", "Copy", "Delete"})

	for _, x := range passwords {
		var row []interface{}

		pwd, err := password.Decrypt([]byte(x.Password), []byte(a.userPassword))
		if err != nil {
			a.errorLog.Fatalln(err)
		}

		row = append(row, strconv.FormatUint(uint64(x.ID), 10))
		row = append(row, x.Item)
		row = append(row, x.Login)
		row = append(row, string(pwd))
		row = append(row, widget.NewButton("Copy", func() {}))
		row = append(row, widget.NewButton("Delete", func() {}))

		slice = append(slice, row)
	}

	return slice
}
