package main

import (
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var vpnCmd *exec.Cmd

func startVPN(status *widget.Label) {
	if vpnCmd != nil {
		status.SetText("Статус: уже запущено")
		return
	}

	// путь к bin/sing-box.exe рядом с portalx.exe
	singboxPath := filepath.Join("bin", "sing-box.exe")

	vpnCmd = exec.Command(singboxPath, "run", "-c", "config.json")

	err := vpnCmd.Start()
	if err != nil {
		status.SetText("Ошибка запуска")
		vpnCmd = nil
		return
	}

	status.SetText("Статус: подключено")
}

func stopVPN(status *widget.Label) {
	if vpnCmd == nil {
		status.SetText("Статус: не запущено")
		return
	}

	err := vpnCmd.Process.Kill()
	if err != nil {
		status.SetText("Ошибка остановки")
		return
	}

	vpnCmd = nil
	status.SetText("Статус: отключено")
}

func main() {

	a := app.New()
	w := a.NewWindow("PortalX")

	linkInput := widget.NewEntry()
	linkInput.SetPlaceHolder("Вставьте VPN ссылку")

	status := widget.NewLabel("Статус: отключено")

	importBtn := widget.NewButton("Импорт", func() {
		status.SetText("Профиль импортирован")
	})

	connectBtn := widget.NewButton("Подключиться", func() {
		startVPN(status)
	})

	disconnectBtn := widget.NewButton("Отключиться", func() {
		stopVPN(status)
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
