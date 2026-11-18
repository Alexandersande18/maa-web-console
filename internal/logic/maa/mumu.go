package maa

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const muMuProfile = "MuMuEmulator12"

func muMuEnabled(cfg *Asst) bool {
	return cfg != nil && cfg.MuMu != nil && cfg.MuMu.Enabled
}

func validateMuMuConfig(cfg *Asst) error {
	if !muMuEnabled(cfg) {
		return nil
	}
	if cfg.ADBType != muMuProfile {
		return fmt.Errorf("mumu.enabled requires adbtype: %s", muMuProfile)
	}
	if strings.TrimSpace(cfg.MuMu.Path) == "" {
		return fmt.Errorf("mumu.path must be the MuMu installation root")
	}
	if cfg.MuMu.Index != nil && *cfg.MuMu.Index < 0 {
		return fmt.Errorf("mumu.index must be non-negative")
	}
	for _, relative := range []string{
		"nx_device/15.0/shell/sdk/external_renderer_ipc.dll",
		"nx_device/12.0/shell/sdk/external_renderer_ipc.dll",
		"shell/sdk/external_renderer_ipc.dll",
	} {
		if st, err := os.Stat(filepath.Join(cfg.MuMu.Path, filepath.FromSlash(relative))); err == nil && !st.IsDir() {
			return nil
		}
	}
	return fmt.Errorf("MuMu external_renderer_ipc.dll not found under mumu.path: %s", cfg.MuMu.Path)
}

// Caller holds nativeMu. Extras are process-global, so set (or clear) them for
// this account after resource reload and immediately before connecting.
func (m *Maa) configureConnectionExtras(profile string) error {
	if muMuEnabled(m.cfg) && profile != muMuProfile {
		return fmt.Errorf("MuMu screenshot enhancement requires connection profile %s", muMuProfile)
	}
	if profile != muMuProfile {
		return nil
	}
	if asstSetConnectionExtras == nil {
		if muMuEnabled(m.cfg) {
			return fmt.Errorf("MaaCore does not export AsstSetConnectionExtras")
		}
		return nil
	}
	values := make(map[string]interface{})
	if muMuEnabled(m.cfg) {
		values["path"] = m.cfg.MuMu.Path
		values["touch"] = false // Screenshot enhancement does not change touch mode.
		client := m.cfg.ClientType
		if client == "" {
			client = "Official"
		}
		values["client_type"] = client
		if m.cfg.MuMu.Index != nil {
			values["index"] = *m.cfg.MuMu.Index
		}
	}
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}
	// The native API returns void; its return register is not a success flag.
	_, err = callStrings(asstSetConnectionExtras, nil, muMuProfile, string(data))
	return err
}
