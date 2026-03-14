package main

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func main() {
    myApp := app.New()
    w := myApp.NewWindow("PortalX GUI")

    w.SetContent(container.NewVBox(
        widget.NewLabel("Добро пожаловать в PortalX!"),
        widget.NewButton("Закрыть", func() {
            myApp.Quit()
        }),
    ))

    w.ShowAndRun()
}