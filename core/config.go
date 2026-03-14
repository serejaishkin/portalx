package core

type RouteRule struct {
	Type string
	Tag  string
}

type Profile struct {
	Name     string
	Type     string
	Route    []RouteRule
	Outbound Outbound
}

type Outbound struct {
	Type       string
	Tag        string
	Server     string
	Port       int
	UUID       string
	PrivateKey string
	PublicKey  string
	Password   string
}

func GenerateConfig(p *Profile) error {

	// позже здесь будет генерация конфигурации sing-box
	return nil
}