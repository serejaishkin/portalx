package core

func GetDefaultRouting() []RouteRule {

	return []RouteRule{
		{
			Type: "direct",
			Tag:  "direct",
		},
		{
			Type: "proxy",
			Tag:  "proxy",
		},
	}
}