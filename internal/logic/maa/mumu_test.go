package maa

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMuMuConnectionExtrasIsolation(t *testing.T) {
	oldExtras, oldConnect := asstSetConnectionExtras, asstConnect
	defer func() { asstSetConnectionExtras, asstConnect = oldExtras, oldConnect }()
	var extras map[string]interface{}
	asstSetConnectionExtras = func(a ...uintptr) (uintptr, uintptr, error) {
		if CharPtr2String(a[0]) != muMuProfile {
			t.Fatal("wrong profile")
		}
		if err := json.Unmarshal([]byte(CharPtr2String(a[1])), &extras); err != nil {
			t.Fatal(err)
		}
		return 0, 0, nil // Native return type is void.
	}
	index := 0
	first := &Maa{cfg: &Asst{MuMu: &MuMuConfig{Enabled: true, Path: `C:\MuMu A`, Index: &index}, ClientType: "Bilibili"}}
	second := &Maa{cfg: &Asst{MuMu: &MuMuConfig{Enabled: true, Path: `C:\MuMu B`}}}
	disabled := &Maa{cfg: &Asst{}}
	for _, tc := range []struct {
		m            *Maa
		path, client string
		index        bool
	}{
		{first, `C:\MuMu A`, "Bilibili", true},
		{second, `C:\MuMu B`, "Official", false},
		{disabled, "", "", false},
		{first, `C:\MuMu A`, "Bilibili", true},
	} {
		asstConnect = func(...uintptr) (uintptr, uintptr, error) {
			if tc.path == "" {
				if len(extras) != 0 {
					t.Fatal("disabled account inherited extras", extras)
				}
			} else {
				if extras["path"] != tc.path || extras["client_type"] != tc.client || extras["touch"] != false {
					t.Fatal("account extras leaked", extras)
				}
				_, hasIndex := extras["index"]
				if hasIndex != tc.index {
					t.Fatal("index omitted or retained incorrectly", extras)
				}
			}
			return 1, 0, nil
		}
		// Clear the decoded map: encoding/json otherwise preserves absent keys.
		extras = nil
		if err := tc.m.AsstConnect("adb", "127.0.0.1:16384", muMuProfile, ""); err != nil {
			t.Fatal(err)
		}
	}
	asstSetConnectionExtras = nil
	if err := first.AsstConnect("adb", "device", muMuProfile, ""); err == nil {
		t.Fatal("missing export accepted")
	}
	if err := first.AsstConnect("adb", "device", "BlueStacks", ""); err == nil {
		t.Fatal("wrong profile accepted")
	}
}

func TestValidateMuMuConfig(t *testing.T) {
	root := t.TempDir()
	cfg := &Asst{ADBType: muMuProfile, MuMu: &MuMuConfig{Enabled: true, Path: root}}
	if err := validateMuMuConfig(cfg); err == nil {
		t.Fatal("missing SDK accepted")
	}
	sdk := filepath.Join(root, "nx_device", "12.0", "shell", "sdk")
	if err := os.MkdirAll(sdk, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sdk, "external_renderer_ipc.dll"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	if err := validateMuMuConfig(cfg); err != nil {
		t.Fatal(err)
	}
	negative := -1
	cfg.MuMu.Index = &negative
	if err := validateMuMuConfig(cfg); err == nil {
		t.Fatal("negative index accepted")
	}
	cfg.MuMu.Enabled = false
	if err := validateMuMuConfig(cfg); err != nil {
		t.Fatal("disabled enhancement rejected", err)
	}
}
