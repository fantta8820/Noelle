package ui

import (
	// "NavigationPlanGenerator/script"
	// "log"
	// "strconv"

	"NavigationPlanGenerator/script"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	"image/color"

	"NavigationPlanGenerator/assets"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	window.SetFixedSize(true)

	title := canvas.NewText("Gerador de Plano de Navegação", color.White)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	input := widget.NewEntry()
	//input.Alignment = fyne.TextAlignCenter
	input.SetPlaceHolder("Insira o nome do plano de navegação")

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
		num, err := strconv.Atoi(input.Text)

		if err != nil {
			log.Fatalf("Erro")
		}

		query, err := script.RunScript(path, num)

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

	contentGrid := container.NewGridWrap(fyne.NewSize(550, 40), input, dialogButton, generateButton, paddingTitle, container.NewGridWrap(fyne.NewSize(550, 300), queryText))

	content := container.NewVBox(paddingTitle, title, paddingTitle, container.NewCenter(contentGrid))

	// text := canvas.NewText("Text object", color.White)
	// text.Alignment = fyne.TextAlignCenter
	// text.TextStyle = fyne.TextStyle{Italic: true}

	window.SetContent(content)

	window.ShowAndRun()
}
