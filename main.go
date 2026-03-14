package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	a := app.New()
	w := a.NewWindow("PortalX")

	linkInput := widget.NewEntry()
	linkInput.SetPlaceHolder("Вставьте VLESS / VPN ссылку")

	status := widget.NewLabel("Статус: отключено")

	importBtn := widget.NewButton("Импорт", func() {
		status.SetText("Профиль импортирован")
	})

	connectBtn := widget.NewButton("Подключиться", func() {
		status.SetText("Подключено")
	})

	disconnectBtn := widget.NewButton("Отключиться", func() {
		status.SetText("Отключено")
	})

	content := container.NewVBox(
		widget.NewLabel("PortalX VPN"),
		linkInput,
		importBtn,
		connectBtn,
		disconnectBtn,
		status,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(420, 300))
	w.ShowAndRun()
}
