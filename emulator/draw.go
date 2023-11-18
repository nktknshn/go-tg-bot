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
	"go.uber.org/zap"
)

type Emulator struct {
	startSent                    bool
	waitingCallback              bool
	handler                      *ActionsHandler
	emulatorStateUpdatedCallback func()
}

func NewEmulator() *Emulator {
	return &Emulator{}
}

// set emulatorStateUpdatedCallback
func (e *Emulator) SetEmulatorStateUpdatedCallback(callback func()) {
	e.emulatorStateUpdatedCallback = callback
}

func (e *Emulator) SetHandler(handler *ActionsHandler) {
	e.handler = handler
}

func (e *Emulator) SetCallbackReceived() {
	e.waitingCallback = false

	if e.emulatorStateUpdatedCallback != nil {
		e.emulatorStateUpdatedCallback()
	}
}

func (e *Emulator) CallbackData(data string) {
	e.waitingCallback = true

	if e.emulatorStateUpdatedCallback != nil {
		e.emulatorStateUpdatedCallback()
	}

	e.handler.CallbackHandlers(data)
}

func (e *Emulator) drawMessageButtons(buttons MessageButtons) []*fyne.Container {
	rows := make([]*fyne.Container, 0)

	for _, r := range buttons.Rows {
		row := container.NewHBox()

		for _, b := range r.Butts {
			b := b
			butt := widget.NewButton(
				b.Title,
				func() {
					data := b.CallbackString

					if data == "" {
						data = b.Title
					}

					e.CallbackData(data)
				},
			)
			row.Add(butt)
		}

		rows = append(rows, row)

	}

	return rows
}

func (e *Emulator) Draw(inp *DrawInput, handler *ActionsHandler) *fyne.Container {

	resultContainer := container.NewVBox()

	for _, el := range inp.Boxes {

		messageBox := container.New(
			layout.NewVBoxLayout(),
		)

		t := widget.NewRichTextWithText(el.Text)

		messageBox.Add(t)

		resultContainer.Add(container.NewStack(
			canvas.NewRectangle(color.RGBA{99, 99, 99, 127}),
			messageBox,
		))

		if len(el.Buttons.Rows) > 0 {
			buttons := e.drawMessageButtons(el.Buttons)

			for _, b := range buttons {
				resultContainer.Add(b)
			}
		}
	}

	textInput := widget.NewEntry()
	textInput.SetPlaceHolder("Write message here...")

	if !e.startSent {
		textInput.SetText("/start")
	}

	sendButton := widget.NewButton(">", func() {
		handler.UserInputHandler(textInput.Text)
		e.startSent = true
	})

	textInputContainer := container.NewBorder(nil, nil, nil, sendButton, textInput)

	resultContainer.Add(textInputContainer)

	logger.Debug("e.waitingCallback",
		zap.Bool("e.waitingCallback", e.waitingCallback),
	)

	if e.waitingCallback {
		resultContainer.Add(
			widget.NewRichTextWithText("Waiting for callback..."),
		)
	}

	return resultContainer
}
