package ui

import (
	"NavigationPlanGenerator/assets"
	"NavigationPlanGenerator/script"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"image/color"	
)

func Run() {
	app := app.New()
	window := app.NewWindow("Noelle")

	window.SetIcon(
		fyne.NewStaticResource(
			"snowflake.png",
			assets.IconData,
		),
	)

	window.Resize(fyne.NewSize(600, 600))

	title := canvas.NewText("Gerador de Plano de Navegação", color.White)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	navigationPlanNameInput := widget.NewEntry()
	navigationPlanNameInput.SetPlaceHolder("Insira o nome do plano de navegação")

	campaignIdInput := widget.NewEntry()
	campaignIdInput.SetPlaceHolder("Insira o ID da campanha")

	campaignIdInput.OnChanged = func(s string) {
		var numbers strings.Builder

		for _, r := range s {
			if unicode.IsDigit(r) {
				numbers.WriteRune(r)
			}
		}

		if numbers.String() != s {
			campaignIdInput.SetText(numbers.String())
		}
	}

	path := ""

	var dialogButton *widget.Button

	dialogButton = widget.NewButton("Selecionar arquivo", func() {
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println("err")
				return
			}

			if reader == nil {
				return
			}

			path = reader.URI().Path()

			dialogButton.SetText(filepath.Base(path))

			reader.Close()
		}, window)

		dialog.SetFilter(
			storage.NewExtensionFileFilter([]string{".xlsx"}),
		)

		dialog.Show()
	})

	queryText := widget.NewMultiLineEntry()
	queryText.SetPlaceHolder("Esperando resultado.")

	generateButton := widget.NewButton("Gerar Plano de Navegação", func() {
		id, err := strconv.Atoi(campaignIdInput.Text)

		query, err := script.RunScript(path, navigationPlanNameInput.Text, id)

		if err != nil {
			queryText.SetText(err.Error())
			return
		}

		queryText.SetText(query)
	})

	generateButton.Importance = widget.HighImportance

	paddingTitle := canvas.NewRectangle(color.Transparent)
	paddingTitle.SetMinSize(fyne.NewSize(0, 20))

	paddingInput := canvas.NewRectangle(color.Transparent)
	paddingInput.SetMinSize(fyne.NewSize(0, 0))

	contentGrid := container.NewGridWrap(fyne.NewSize(550, 40), navigationPlanNameInput, campaignIdInput, dialogButton, generateButton, paddingTitle, container.NewGridWrap(fyne.NewSize(550, 250), queryText))

	content := container.NewVBox(paddingTitle, title, paddingTitle, container.NewCenter(contentGrid))

	window.SetContent(content)

	window.ShowAndRun()
}
