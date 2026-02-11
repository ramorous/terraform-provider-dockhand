package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildClientConfigFromModel(t *testing.T) {
	m := DockhandProviderModel{
		Endpoint: types.StringValue("http://localhost:3000"),
		Cookie:   types.StringValue("session=abc123"),
		Timeout:  types.Int64Value(15),
	}

	cfg := BuildClientConfigFromModel(m)
	if cfg.Endpoint != "http://localhost:3000" {
		t.Fatalf("unexpected endpoint: %s", cfg.Endpoint)
	}
	if cfg.Cookie != "session=abc123" {
		t.Fatalf("unexpected cookie: %s", cfg.Cookie)
	}
	if cfg.Timeout != 15 {
		t.Fatalf("unexpected timeout: %d", cfg.Timeout)
	}
}

func TestBuildClientConfigFromModelEmpty(t *testing.T) {
	m := DockhandProviderModel{}
	cfg := BuildClientConfigFromModel(m)
	if cfg.Endpoint != "" || cfg.Cookie != "" || cfg.Timeout != 0 {
		t.Fatalf("expected empty client config for zero model, got: %+v", cfg)
	}
}
