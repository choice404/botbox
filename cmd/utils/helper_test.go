/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// touch creates an empty file so the detector has something to find
func touch(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create %s: %v", path, err)
	}
}

func TestDetectEnvProvider(t *testing.T) {
	t.Run("doppler.yaml means doppler", func(t *testing.T) {
		dir := t.TempDir()
		touch(t, filepath.Join(dir, "doppler.yaml"))
		if got := DetectEnvProvider(dir); got != "doppler" {
			t.Errorf("DetectEnvProvider() = %q, want doppler", got)
		}
	})

	t.Run(".env means env", func(t *testing.T) {
		dir := t.TempDir()
		touch(t, filepath.Join(dir, ".env"))
		if got := DetectEnvProvider(dir); got != "env" {
			t.Errorf("DetectEnvProvider() = %q, want env", got)
		}
	})

	t.Run("doppler.yaml wins over .env", func(t *testing.T) {
		dir := t.TempDir()
		touch(t, filepath.Join(dir, "doppler.yaml"))
		touch(t, filepath.Join(dir, ".env"))
		if got := DetectEnvProvider(dir); got != "doppler" {
			t.Errorf("DetectEnvProvider() = %q, want doppler", got)
		}
	})

	t.Run("neither file defaults to env", func(t *testing.T) {
		dir := t.TempDir()
		if got := DetectEnvProvider(dir); got != "env" {
			t.Errorf("DetectEnvProvider() = %q, want env", got)
		}
	})
}

func TestResolveEnvProvider(t *testing.T) {
	t.Run("configured value wins over detection", func(t *testing.T) {
		dir := t.TempDir()
		touch(t, filepath.Join(dir, ".env"))
		config := Config{BotInfo: BotConfig{EnvProvider: "doppler"}}
		if got := ResolveEnvProvider(config, dir); got != "doppler" {
			t.Errorf("ResolveEnvProvider() = %q, want doppler", got)
		}
	})

	t.Run("empty value falls back to detection", func(t *testing.T) {
		dir := t.TempDir()
		touch(t, filepath.Join(dir, "doppler.yaml"))
		config := Config{}
		if got := ResolveEnvProvider(config, dir); got != "doppler" {
			t.Errorf("ResolveEnvProvider() = %q, want doppler", got)
		}
	})

	t.Run("empty value and no files defaults to env", func(t *testing.T) {
		dir := t.TempDir()
		config := Config{}
		if got := ResolveEnvProvider(config, dir); got != "env" {
			t.Errorf("ResolveEnvProvider() = %q, want env", got)
		}
	})
}

/*
Copyright © 2025 Austin "Choice404" Choi

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
