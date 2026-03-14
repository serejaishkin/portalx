package core

type Profile struct {
	Name    string
	Type    string
	Route   string
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

func GenerateConfig(p Profile) error {

	// позже здесь будет генерация sing-box config
	return nil
}