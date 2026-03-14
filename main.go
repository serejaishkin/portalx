package main

import (
    "log"
    "portalx/core"
)

func main() {
    log.Println("PortalX starting...")

    profile := core.LoadProfile("profiles/default.yaml")

    err := core.GenerateConfig(profile)
    if err != nil {
        log.Fatal("Failed to generate config:", err)
    }

    err = core.StartSingBox()
    if err != nil {
        log.Fatal("Failed to start sing-box:", err)
    }

    log.Println("PortalX running...")
    select {} // держим процесс живым
}