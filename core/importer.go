package core

import (
    "bufio"
    "errors"
    "log"
    "net/url"
    "os"
    "strings"
)

// ---------------------------
// Импорт Amnezia 1.5 .conf
// ---------------------------
func ImportAmneziaConf(filePath string) (*Profile, error) {
    f, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    var ifacePrivate, ifaceId string
    var peerPublic, peerEndpoint string
    var peerAllowed []string

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
            continue
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(parts[0])
        val := strings.TrimSpace(parts[1])

        switch key {
        case "PrivateKey":
            ifacePrivate = val
        case "Id":
            ifaceId = val
        case "PublicKey":
            peerPublic = val
        case "AllowedIPs":
            peerAllowed = strings.Split(val, ",")
            for i := range peerAllowed {
                peerAllowed[i] = strings.TrimSpace(peerAllowed[i])
            }
        case "Endpoint":
            peerEndpoint = val
        }
    }

    if ifacePrivate == "" || peerPublic == "" || peerEndpoint == "" {
        return nil, errors.New("invalid Amnezia config: missing required fields")
    }

    profile := &Profile{
        Name: ifaceId,
        Outbounds: []Outbound{
            {
                Type:       "wireguard",
                Tag:        "wg",
                Server:     peerEndpoint,
                PrivateKey: ifacePrivate,
                PublicKey:  peerPublic,
                AllowedIPs: peerAllowed,
            },
        },
        Route: GetDefaultRouting(),
    }

    return profile, nil
}

// ---------------------------
// Импорт VLESS ссылки
// ---------------------------
func ImportVLESS(link string) (*Profile, error) {
    u, err := url.Parse(link)
    if err != nil {
        return nil, err
    }

    uuid := ""
    if u.User != nil {
        uuid = u.User.String()
    }

    tag := u.Fragment
    if tag == "" {
        tag = "vless"
    }

    profile := &Profile{
        Name: tag,
        Outbounds: []Outbound{
            {
                Type:   "vless",
                Tag:    tag,
                Server: u.Host,
                UUID:   uuid,
            },
        },
        Route: GetDefaultRouting(),
    }

    return profile, nil
}

// ---------------------------
// Пример вызова из main.go
// ---------------------------
func ExampleUsage() {
    // Импорт Amnezia
    profile1, err := ImportAmneziaConf("profiles/WARPw_96.conf")
    if err != nil {
        log.Println("Amnezia import error:", err)
    } else {
        GenerateConfig(profile1)
    }

    // Импорт VLESS
    profile2, err := ImportVLESS("vless://UUID@server:port?encryption=none&type=ws#myServer")
    if err != nil {
        log.Println("VLESS import error:", err)
    } else {
        GenerateConfig(profile2)
    }
}