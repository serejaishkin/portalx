package core

import "os/exec"

func StartSingBox() error {
    cmd := exec.Command("bin/sing-box.exe", "run", "-c", "config.json")
    return cmd.Start()
}