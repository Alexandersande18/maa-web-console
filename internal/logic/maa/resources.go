package maa

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

var runtimeRoot, runtimeUserDir string // guarded by nativeMu

func initializeRuntime(root, userDir, client string) error {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	if runtimeRoot != "" {
		if !strings.EqualFold(runtimeRoot, root) {
			return fmt.Errorf("all instances must use the same resource root: %s != %s", runtimeRoot, root)
		}
		if userDir != "" {
			userDir, err = filepath.Abs(userDir)
			if err != nil {
				return err
			}
			if !strings.EqualFold(runtimeUserDir, userDir) {
				return fmt.Errorf("MaaCore user directory is process-global; use maa.userdir")
			}
		}
		return nil
	}
	if userDir == "" {
		userDir = filepath.Join(".runtime", "maa")
	}
	userDir, err = filepath.Abs(userDir)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(userDir, 0755); err != nil {
		return err
	}
	if err = newSyscall(root); err != nil {
		return err
	}
	r, err := callStrings(asstSetUserDir, nil, userDir)
	if err != nil {
		return err
	}
	if r == 0 {
		return fmt.Errorf("AsstSetUserDir failed: %s", userDir)
	}
	if err = loadResources(root, client); err != nil {
		return err
	}
	runtimeRoot, runtimeUserDir = root, userDir
	log.WithFields(log.Fields{"resource": root, "userdir": userDir}).Info("MAA runtime initialized")
	return nil
}

func loadResource(root string) error {
	r, err := callStrings(asstLoadResource, nil, root)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"path": root, "success": r != 0}).Info("AsstLoadResource")
	if r == 0 {
		return fmt.Errorf("%w: %s", ErrorAsstLoadSource, root)
	}
	return nil
}

func prepareHotUpdate(root string) error {
	src := filepath.Join(root, "resource", "tasks.json")
	data, err := os.ReadFile(src)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	dir := filepath.Join(root, "resource", "tasks")
	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	dst := filepath.Join(dir, "tasks.json")
	if previous, readErr := os.ReadFile(dst); readErr == nil && bytes.Equal(previous, data) {
		return nil
	}
	return os.WriteFile(dst, data, 0644)
}

// Paths are the parents of resource/. Match the official CN loading order:
// base -> cache; overseas resources follow their corresponding base layer.
func loadResources(root, client string) error {
	roots := []string{root, filepath.Join(root, "cache")}
	if client != "" && client != "Official" && client != "Bilibili" {
		roots = append(roots, filepath.Join(root, "resource", "global", client), filepath.Join(root, "cache", "resource", "global", client))
	}
	for i, path := range roots {
		st, err := os.Stat(filepath.Join(path, "resource"))
		if i > 0 && os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if !st.IsDir() {
			return fmt.Errorf("resource is not a directory: %s", path)
		}
		// Only API cache files need the old-to-new tasks layout migration.
		if i == 1 || i == 3 {
			if err = prepareHotUpdate(path); err != nil {
				return err
			}
		}
		if err = loadResource(path); err != nil {
			return err
		}
	}
	return nil
}
func reloadResources(ctx context.Context, root, client string) error {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	return loadResources(root, client)
}
