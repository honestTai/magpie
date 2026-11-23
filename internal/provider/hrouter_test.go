package provider

import "testing"

// HRouter is a pay-as-you-go Sub2API relay at hrouter.net: Chat and Responses
// at /v1, Anthropic Messages at the root, and CNY key balances at /v1/usage.
func TestHRouterPreset(t *testing.T) {
	p, err := FromPreset("hrouter")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "HRouter" || p.Icon != "hrouter-color" ||
		p.Chat != "https://hrouter.net/v1" ||
		p.Responses != "https://hrouter.net/v1" ||
		p.Anthropic != "https://hrouter.net" ||
		p.Website != "https://hrouter.net" ||
		p.KeysURL != "https://hrouter.net/keys" {
		t.Fatalf("%+v", p)
	}
	pr := Preset("hrouter")
	if pr == nil || pr.Kind != KindRelay || pr.NoKey || pr.Sponsored {
		t.Fatalf("%+v", pr)
	}
	p.Key = "sk-test"
	if got := p.BillingKind(); got != BillingPAYG {
		t.Fatalf("billing = %q, want %q", got, BillingPAYG)
	}
	src, ok := balanceSourceOf(p)
	if !ok || src.url != "https://hrouter.net/v1/usage" {
		t.Fatalf("balanceSourceOf = %q, %v", src.url, ok)
	}
	amt, err := src.read([]byte(`{"mode":"unrestricted","isValid":true,"remaining":42.5,"unit":"CNY"}`))
	if err != nil || amt != "¥42.50" {
		t.Fatalf("balance read = %q, %v", amt, err)
	}
	im, why := imported("x", "sk-test", endpoints{chat: "https://hrouter.net/v1/"}, nil)
	if why != "" || im.Preset != "hrouter" {
		t.Fatalf("imported: %+v (%s)", im, why)
	}
}
