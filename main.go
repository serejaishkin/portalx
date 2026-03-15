package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PortalXRules struct {
	VPNDomains    []string `json:"vpn_domains"`
	ZapretDomains []string `json:"zapret_domains"`
	DirectDomains []string `json:"direct_domains"`
}

var vpnCmd *exec.Cmd
var zapretCmd *exec.Cmd
var profiles []string
var rules PortalXRules

func ensureDirs() {

	os.MkdirAll("profiles", 0755)
	os.MkdirAll(filepath.Join("tools", "zapret"), 0755)
}

func loadProfiles() {

	files, err := os.ReadDir("profiles")
	if err != nil {
		return
	}

	profiles = []string{}

	for _, f := range files {

		if strings.HasSuffix(f.Name(), ".json") {
			profiles = append(profiles, f.Name())
		}

	}
}

func loadRules() {

	path := "portalx_rules.json"

	data, err := os.ReadFile(path)

	if err != nil {

		rules = PortalXRules{
			VPNDomains: []string{"openai.com", "claude.ai"},
			ZapretDomains: []string{
				"youtube.com",
				"googlevideo.com",
				"discord.com",
			},
		}

		saveRules()

		return
	}

	json.Unmarshal(data, &rules)
}

func saveRules() {

	data, _ := json.MarshalIndent(rules, "", "  ")

	os.WriteFile("portalx_rules.json", data, 0644)
}

func generateZapretHostlist() {

	path := filepath.Join("tools", "zapret", "hostlist.txt")

	content := strings.Join(rules.ZapretDomains, "\n")

	os.WriteFile(path, []byte(content), 0644)
}

func createConfig(link string) (string, error) {

	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	uuid := u.User.Username()
	server := u.Hostname()

	portStr := u.Port()
	port, _ := strconv.Atoi(portStr)

	query := u.Query()

	security := query.Get("security")
	sni := query.Get("sni")

	config := map[string]interface{}{

		"log": map[string]interface{}{
			"level": "info",
		},

		"inbounds": []map[string]interface{}{
			{
				"type": "tun",
				"tag":  "tun-in",
				"interface_name": "portalx",
				"inet4_address":  "172.19.0.1/30",
				"auto_route":     true,
			},
		},
	}

	outbound := map[string]interface{}{
		"type":        "vless",
		"tag":         "proxy",
		"server":      server,
		"server_port": port,
		"uuid":        uuid,
	}

	if security == "tls" {

		outbound["tls"] = map[string]interface{}{
			"enabled":     true,
			"server_name": sni,
		}

	}

	config["outbounds"] = []map[string]interface{}{outbound}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("profile_%s.json", uuid[:8])

	path := filepath.Join("profiles", name)

	err = os.WriteFile(path, data, 0644)

	if err != nil {
		return "", err
	}

	return name, nil
}

func startVPN(profile string, status *widget.Label) {

	if vpnCmd != nil {
		status.SetText("VPN уже запущен")
		return
	}

	configPath := filepath.Join("profiles", profile)

	singbox := filepath.Join("bin", "sing-box.exe")

	vpnCmd = exec.Command(singbox, "run", "-c", configPath)

	err := vpnCmd.Start()

	if err != nil {

		status.SetText("Ошибка запуска VPN")

		vpnCmd = nil

		return
	}

	status.SetText("VPN подключен")
}

func stopVPN(status *widget.Label) {

	if vpnCmd == nil {

		status.SetText("VPN не запущен")

		return
	}

	vpnCmd.Process.Kill()

	vpnCmd = nil

	status.SetText("VPN отключен")
}

func startZapret(status *widget.Label) {

	if zapretCmd != nil {
		status.SetText("zapret уже запущен")
		return
	}

	generateZapretHostlist()

	bat := filepath.Join("tools", "zapret", "general.bat")

	if _, err := os.Stat(bat); os.IsNotExist(err) {

		status.SetText("zapret не найден")

		return
	}

	zapretCmd = exec.Command("cmd", "/C", bat)

	err := zapretCmd.Start()

	if err != nil {

		status.SetText("Ошибка запуска zapret")

		zapretCmd = nil

		return
	}

	status.SetText("zapret запущен")
}

func stopZapret(status *widget.Label) {

	if zapretCmd == nil {

		status.SetText("zapret не запущен")

		return
	}

	zapretCmd.Process.Kill()

	zapretCmd = nil

	status.SetText("zapret остановлен")
}

func main() {

	ensureDirs()
	loadProfiles()
	loadRules()

	a := app.New()

	w := a.NewWindow("PortalX")

	linkInput := widget.NewEntry()
	linkInput.SetPlaceHolder("Вставьте VLESS ссылку")

	status := widget.NewLabel("Статус: отключено")

	profileList := widget.NewList(

		func() int {
			return len(profiles)
		},

		func() fyne.CanvasObject {
			return widget.NewLabel("profile")
		},

		func(i widget.ListItemID, o fyne.CanvasObject) {

			o.(*widget.Label).SetText(profiles[i])

		},
	)

	selectedProfile := ""

	profileList.OnSelected = func(id int) {

		selectedProfile = profiles[id]

	}

	importBtn := widget.NewButton("Импорт", func() {

		link := strings.TrimSpace(linkInput.Text)

		if !strings.HasPrefix(link, "vless://") {

			status.SetText("Это не VLESS ссылка")

			return
		}

		name, err := createConfig(link)

		if err != nil {

			status.SetText("Ошибка: " + err.Error())

			return
		}

		loadProfiles()

		profileList.Refresh()

		status.SetText("Профиль создан: " + name)

	})

	connectBtn := widget.NewButton("Подключить VPN", func() {

		if selectedProfile == "" {

			status.SetText("Выберите профиль")

			return
		}

		startVPN(selectedProfile, status)

	})

	stopVPNBtn := widget.NewButton("Отключить VPN", func() {

		stopVPN(status)

	})

	startZapretBtn := widget.NewButton("Запустить zapret", func() {

		startZapret(status)

	})

	stopZapretBtn := widget.NewButton("Остановить zapret", func() {

		stopZapret(status)

	})

	content := container.NewVBox(

		widget.NewLabel("PortalX"),

		linkInput,

		importBtn,

		widget.NewLabel("Профили"),

		profileList,

		connectBtn,
		stopVPNBtn,

		widget.NewLabel("DPI обход"),

		startZapretBtn,
		stopZapretBtn,

		status,
	)

	w.SetContent(content)

	w.Resize(fyne.NewSize(460, 520))

	w.ShowAndRun()
}