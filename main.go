package main

import (
    "portalx/core"
    "portalx/profiles"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
)

func main() {
    a := app.New()
    w := a.NewWindow("PortalX GUI")

    importAmneziaBtn := widget.NewButton("Import Amnezia Config", func() {
        dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err != nil || reader == nil {
                return
            }
            profile, err := core.ImportAmneziaConf(reader.URI().Path())
            if err != nil {
                log.Println("Import error:", err)
                return
            }
            core.GenerateConfig(profile)
            log.Println("Amnezia config imported successfully")
        }, w)
    })

    importVlessBtn := widget.NewButton("Add VLESS Link", func() {
        dialog.ShowEntryDialog("Paste VLESS Link", "VLESS Link:", func(link string) {
            profile, err := core.ImportVLESS(link)
            if err != nil {
                log.Println("VLESS import error:", err)
                return
            }
            core.GenerateConfig(profile)
            log.Println("VLESS link added successfully")
        }, w)
    })

    w.SetContent(container.NewVBox(
        importAmneziaBtn,
        importVlessBtn,
    ))

    w.ShowAndRun()
}