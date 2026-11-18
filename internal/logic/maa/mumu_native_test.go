package maa

import (
	"bytes"
	"context"
	"encoding/json"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcfg"
)

type muMuProbeHook struct{ fastest chan muMuProbeEvent }
type muMuProbeEvent struct {
	What    string `json:"what"`
	UUID    string `json:"uuid"`
	Details struct {
		Method string `json:"method"`
		Cost   int    `json:"cost"`
	} `json:"details"`
}

func (h *muMuProbeHook) Do(_ int, details string, _ interface{}) {
	var event muMuProbeEvent
	if json.Unmarshal([]byte(details), &event) == nil && event.What == "FastestWayToScreencap" {
		select {
		case h.fastest <- event:
		default:
		}
	}
}

// Opt-in read-only device test: connects and captures images without starting tasks.
func TestMuMuNativeScreenshots(t *testing.T) {
	if os.Getenv("MAA_MUMU_SMOKE") != "1" {
		t.Skip("requires MAA_MUMU_SMOKE=1 and idle MuMu devices")
	}
	adapter, err := gcfg.NewAdapterFile(os.Getenv("MAA_MUMU_SMOKE_CONFIG"))
	if err != nil {
		t.Fatal(err)
	}
	value, err := gcfg.NewWithAdapter(adapter).Get(context.Background(), "maa")
	if err != nil {
		t.Fatal(err)
	}
	var cfg Config
	if err = value.Scan(&cfg); err != nil {
		t.Fatal(err)
	}
	if len(cfg.Asst) == 0 {
		t.Fatal("no MuMu instances configured")
	}
	output := os.Getenv("MAA_MUMU_SMOKE_OUTPUT")
	if !filepath.IsAbs(output) {
		t.Fatal("absolute MAA_MUMU_SMOKE_OUTPUT required")
	}
	if err = os.MkdirAll(output, 0755); err != nil {
		t.Fatal(err)
	}
	type result struct {
		Instance   string  `json:"instance"`
		Address    string  `json:"address"`
		Enhanced   bool    `json:"enhanced"`
		Method     string  `json:"method"`
		CoreCostMS int     `json:"coreCostMs"`
		AverageMS  float64 `json:"averageMs"`
		Width      int     `json:"width"`
		Height     int     `json:"height"`
		Samples    int     `json:"samples"`
	}
	var results []result
	seen := make(map[string]string)
	for _, enhanced := range []bool{false, true} {
		for i, original := range cfg.Asst {
			c := *original
			if !muMuEnabled(&c) {
				t.Fatalf("instance %s must configure MuMu", c.ID)
			}
			c.UserDir = filepath.Join(output, "native")
			if c.ClientType == "" {
				c.ClientType = "Official"
			}
			if !enhanced {
				c.MuMu = nil
			}
			func() {
				m, err := New(&c)
				if err != nil {
					t.Fatal(err)
				}
				defer m.AsstDestroy()
				hook := &muMuProbeHook{fastest: make(chan muMuProbeEvent, 1)}
				m.mu.Lock()
				m.hooks = []Hook{hook}
				m.mu.Unlock()
				if err := reloadResources(context.Background(), c.ResourcePath, c.ClientType); err != nil {
					t.Fatal(err)
				}
				if err := m.AsstConnect(c.ADBPath, c.ADBAddress, c.ADBType, ""); err != nil {
					t.Fatal(err)
				}
				var event muMuProbeEvent
				select {
				case event = <-hook.fastest:
				case <-time.After(5 * time.Second):
					t.Fatal("missing screenshot method callback")
				}
				if enhanced && event.Details.Method != "MumuExtras" {
					t.Fatalf("enhancement not active on %s: %s", c.ID, event.Details.Method)
				}
				if !enhanced && event.Details.Method == "MumuExtras" {
					t.Fatal("disabled account inherited enhancement")
				}
				if event.UUID == "" {
					t.Fatal("missing device UUID")
				}
				for address, uuid := range seen {
					if address != c.ADBAddress && uuid == event.UUID {
						t.Fatal("instances connected to the same device")
					}
				}
				seen[c.ADBAddress] = event.UUID
				r := result{Instance: c.ID, Address: c.ADBAddress, Enhanced: enhanced, Method: event.Details.Method, CoreCostMS: event.Details.Cost, Samples: 10}
				var elapsed time.Duration
				for sample := 0; sample < r.Samples; sample++ {
					start := time.Now()
					data, err := m.GetImage()
					elapsed += time.Since(start)
					if err != nil {
						t.Fatal(err)
					}
					image, err := png.DecodeConfig(bytes.NewReader(data))
					if err != nil || image.Width == 0 || image.Height == 0 {
						t.Fatal("invalid screenshot", err)
					}
					r.Width, r.Height = image.Width, image.Height
					if enhanced && sample == 0 {
						if err = os.WriteFile(filepath.Join(output, "instance-"+strconv.Itoa(i)+".png"), data, 0600); err != nil {
							t.Fatal(err)
						}
					}
				}
				r.AverageMS = float64(elapsed.Microseconds()) / 1000 / float64(r.Samples)
				results = append(results, r)
				t.Logf("%s enhanced=%v method=%s core=%dms screenshot=%.2fms %dx%d samples=%d", c.ID, enhanced, r.Method, r.CoreCostMS, r.AverageMS, r.Width, r.Height, r.Samples)
			}()
		}
	}
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(output, "screenshots.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}
