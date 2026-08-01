/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func CompareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := max(len(parts1), len(parts2))

	for len(parts1) < maxLen {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < maxLen {
		parts2 = append(parts2, "0")
	}

	for i := range maxLen {
		var num1, num2 int
		fmt.Sscanf(parts1[i], "%d", &num1)
		fmt.Sscanf(parts2[i], "%d", &num2)

		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
	}

	return 0
}

func GetLatestVersion() (*GitHubRelease, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/choice404/botbox/releases/latest", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "botbox-cli")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &release, nil
}

func UpdateBotBox(version string) error {
	fmt.Printf("🔄 Updating botbox to %s...\n", version)

	cleanCmd := exec.Command("go", "clean", "-modcache")
	if err := cleanCmd.Run(); err != nil {
		fmt.Printf("⚠️  Warning: failed to clean module cache: %v\n", err)
	}

	installCmd := exec.Command("go", "install", fmt.Sprintf("github.com/choice404/botbox/v2@%s", version))
	if output, err := installCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install update: %w\n%s", err, string(output))
	}

	if err := verifyInstalledVersion(version); err != nil {
		fmt.Printf("⚠️  Warning: %v\n", err)
	}

	if err := updateGlobalConfigVersion(version); err != nil {
		fmt.Printf("⚠️  Warning: failed to update global config version: %v\n", err)
	}

	fmt.Printf("✅ Successfully updated to botbox %s\n", version)
	return nil
}

/**
 * installBinDir
 * Finds the directory go install writes binaries to
 * @return string - the bin directory
 * @return error - any failure running go env
 **/
func installBinDir() (string, error) {
	out, err := exec.Command("go", "env", "GOBIN").Output()
	if err != nil {
		return "", fmt.Errorf("failed to run go env: %w", err)
	}
	gobin := strings.TrimSpace(string(out))
	if gobin != "" {
		return gobin, nil
	}

	out, err = exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return "", fmt.Errorf("failed to run go env: %w", err)
	}
	gopath := strings.TrimSpace(string(out))
	if gopath == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}
	return filepath.Join(gopath, "bin"), nil
}

/**
 * verifyInstalledVersion
 * Runs the freshly installed binary and checks it reports the expected version
 * Also warns when the currently running binary lives somewhere else in PATH
 * @param version {string} - the version that was just installed
 * @return error - a description of any mismatch found
 **/
func verifyInstalledVersion(version string) error {
	binDir, err := installBinDir()
	if err != nil {
		return err
	}

	binName := "botbox"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(binDir, binName)

	out, err := exec.Command(binPath, "--version").Output()
	if err != nil {
		return fmt.Errorf("could not run the installed binary at %s: %w", binPath, err)
	}

	installed := strings.TrimSpace(string(out))
	if strings.TrimPrefix(installed, "v") != strings.TrimPrefix(version, "v") {
		return fmt.Errorf("the installed binary reports %s instead of %s. The Go module proxy may be serving an older snapshot of this tag, this happens when a tag is moved after it was already fetched. A new release with a fresh tag is needed", installed, version)
	}

	// Warn when the running binary is not the one that was just updated
	if exe, err := os.Executable(); err == nil {
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			installedResolved, _ := filepath.EvalSymlinks(binPath)
			if installedResolved != "" && resolved != installedResolved {
				return fmt.Errorf("the update was installed to %s but the botbox you are running lives at %s. Check your PATH order or remove the old copy", binPath, resolved)
			}
		}
	}

	return nil
}

func updateGlobalConfigVersion(version string) error {
	exists, err := GlobalConfigExists()
	if err != nil {
		return fmt.Errorf("failed to check global config: %w", err)
	}

	configVersion := strings.TrimPrefix(version, "v")

	if !exists {
		if err := CreateGlobalConfig(); err != nil {
			return fmt.Errorf("failed to create global config: %w", err)
		}
		fmt.Printf("📝 Created global config with version %s\n", configVersion)
		return nil
	}

	if err := SetGlobalConfigValue("cli.version", configVersion); err != nil {
		return fmt.Errorf("failed to set version in global config: %w", err)
	}

	fmt.Printf("📝 Updated global config version to %s\n", configVersion)
	return nil
}

func CheckForUpdates() (*GitHubRelease, bool, error) {
	latest, err := GetLatestVersion()
	if err != nil {
		return nil, false, fmt.Errorf("failed to check for updates: %w", err)
	}

	currentVersion := strings.TrimPrefix(Version, "v")
	latestVersion := strings.TrimPrefix(latest.TagName, "v")

	hasUpdate := CompareVersions(currentVersion, latestVersion) < 0

	return latest, hasUpdate, nil
}

func ShouldCheckForUpdates() bool {
	exists, err := GlobalConfigExists()
	if err != nil || !exists {
		return false
	}

	checkUpdates := GetGlobalConfigValue("cli.check_updates")
	if checkUpdates == nil {
		return false
	}

	if enabled, ok := checkUpdates.(bool); ok {
		return enabled
	}

	return false
}

func ShouldAutoUpdate() bool {
	exists, err := GlobalConfigExists()
	if err != nil || !exists {
		return false
	}

	autoUpdate := GetGlobalConfigValue("cli.auto_update")
	if autoUpdate == nil {
		return false
	}

	if enabled, ok := autoUpdate.(bool); ok {
		return enabled
	}

	return false
}

/*
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.

MIT License

Copyright (c) 2025 Austin Choi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
