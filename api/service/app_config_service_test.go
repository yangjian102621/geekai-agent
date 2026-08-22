package service

import (
	"geekai/core/types"
	"geekai/store/vo"
	"strings"
	"testing"
)

func TestAppConfigServiceEncryptsAndDecrypts(t *testing.T) {
	service, err := NewAppConfigService(&types.AppConfig{AppConfigKey: "local-test-app-config-key"})
	if err != nil {
		t.Fatal(err)
	}

	original := vo.AppConfig{
		ApiUrl:        "https://api.example.com",
		Token:         "token-secret",
		PrivateKey:    "private-key-secret",
		BailianApiKey: "bailian-secret",
	}
	encrypted, err := service.Encode(original)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(encrypted, original.Token) || strings.Contains(encrypted, original.PrivateKey) || strings.Contains(encrypted, original.BailianApiKey) {
		t.Fatal("encrypted app config contains plaintext secret")
	}

	var decoded vo.AppConfig
	if err := service.Decode(encrypted, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded != original {
		t.Fatalf("decoded config does not match original: got %+v, want %+v", decoded, original)
	}
}

func TestAppConfigServiceMaskAndMergeSecrets(t *testing.T) {
	service, err := NewAppConfigService(&types.AppConfig{AppConfigKey: "local-test-app-config-key"})
	if err != nil {
		t.Fatal(err)
	}

	current := vo.AppConfig{Token: "old-token", PrivateKey: "old-private-key", BailianApiKey: "old-bailian-key"}
	masked := service.Mask(current)
	if masked.Token != "" || masked.PrivateKey != "" || masked.BailianApiKey != "" {
		t.Fatalf("mask did not clear secrets: %+v", masked)
	}

	merged := service.MergeSecrets(current, vo.AppConfig{ApiUrl: "https://new.example.com"})
	if merged.Token != current.Token || merged.PrivateKey != current.PrivateKey || merged.BailianApiKey != current.BailianApiKey {
		t.Fatalf("merge did not preserve secrets: %+v", merged)
	}
}
