package core

type Profile struct {
	Name     string
	Type     string
	Server   string
	Port     int
	UUID     string
	Password string
}

type Outbound struct {
	Type   string
	Server string
	Port   int
}

func GenerateConfig(p Profile) error {

	// здесь позже будем генерировать sing-box config
	return nil
}