package core

import "fmt"

var currentProfile *Profile

func ImportLink(link string) error {

	fmt.Println("Import link:", link)

	currentProfile = &Profile{
		Name: "default",
		Type: "vless",
		Route: GetDefaultRouting(),
		Outbound: Outbound{
			Type: "vless",
			Server: "example.com",
			Port: 443,
		},
	}

	return GenerateConfig(currentProfile)
}