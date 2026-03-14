package main

import (
	"portalx/core"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	a := app.New()
	w := a.NewWindow("PortalX")

	input := widget.NewEntry()
	input.SetPlaceHolder("Paste VLESS or WG link")

	status := widget.NewLabel("Disconnected")

	importBtn := widget.NewButton("Import", func() {

		err := core.ImportLink(input.Text)

		if err != nil {
			status.SetText("Import error")
			return
		}

		status.SetText("Profile imported")
	})

	connectBtn := widget.NewButton("Connect", func() {

		err := core.StartVPN()

		if err != nil {
			status.SetText("Connect error")
			return
		}

		status.SetText("Connected")
	})

	disconnectBtn := widget.NewButton("Disconnect", func() {

		err := core.StopVPN()

		if err != nil {
			status.SetText("Disconnect error")
			return
		}

		status.SetText("Disconnected")
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("PortalX VPN"),
		input,
		importBtn,
		connectBtn,
		disconnectBtn,
		status,
	))

	w.Resize(fyne.NewSize(400, 260))
	w.ShowAndRun()
}