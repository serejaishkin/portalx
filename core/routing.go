package core

type RouteRule struct {
    DomainSuffix []string `json:"domain_suffix,omitempty"`
    Outbound     string   `json:"outbound"`
}

func GetDefaultRouting() []RouteRule {
    return []RouteRule{
        {DomainSuffix: []string{"openai.com", "discord.com"}, Outbound: "wg"},
        {DomainSuffix: []string{"googlevideo.com"}, Outbound: "direct"},
    }
}