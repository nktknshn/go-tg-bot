package emulator

import (

	// "log"
	// "image/color"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// func drawMessageBox(box MessageBox) {

// }

func EmulatorDraw(inp *DrawInput, handler *ActionsHandler) *fyne.Container {

	resultContainer := container.NewVBox()

	for _, el := range inp.Boxes {

		messageBox := container.New(
			layout.NewVBoxLayout(),
		)

		t := widget.NewRichTextWithText(el.Text)

		messageBox.Add(t)

		for _, r := range el.Buttons.Rows {
			row := container.NewHBox()

			for _, b := range r.Butts {
				butt := widget.NewButton(
					b.Title, func() { handler.CallbackHandlers(b.CallbackString) },
				)
				row.Add(butt)
			}

			messageBox.Add(row)
		}

		resultContainer.Add(container.NewStack(
			canvas.NewRectangle(color.RGBA{127, 127, 127, 127}),
			messageBox,
		))
	}

	textInput := widget.NewEntry()
	textInput.SetPlaceHolder("Write message here...")
	textInput.SetText("/start")

	sendButton := widget.NewButton(">", func() {
		handler.UserInputHandler(textInput.Text)
	})

	textInputContainer := container.NewBorder(nil, nil, nil, sendButton, textInput)

	// textInputContainer := container.NewAdaptiveGrid(
	// 	2,
	// 	textInput,
	// 	// layout.NewSpacer(),
	// 	sendButton,
	// )

	resultContainer.Add(textInputContainer)

	return resultContainer
}
