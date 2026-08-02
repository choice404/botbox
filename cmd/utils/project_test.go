/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// readOutput loads a generated file so the tests can assert on its content
func readOutput(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return string(content)
}

func TestGenerateDockerFilesEnv(t *testing.T) {
	dir := t.TempDir()

	written, err := GenerateDockerFiles(dir, "3.12", "env", false)
	if err != nil {
		t.Fatalf("GenerateDockerFiles() error = %v", err)
	}
	if len(written) != 3 {
		t.Fatalf("GenerateDockerFiles() wrote %d files, want 3", len(written))
	}

	dockerfile := readOutput(t, filepath.Join(dir, "Dockerfile"))
	if !strings.Contains(dockerfile, "FROM python:3.12-slim") {
		t.Errorf("Dockerfile missing python version base image:\n%s", dockerfile)
	}
	if !strings.Contains(dockerfile, "WORKDIR /app") {
		t.Errorf("Dockerfile missing WORKDIR /app:\n%s", dockerfile)
	}
	if !strings.Contains(dockerfile, "COPY botbox.conf .") {
		t.Errorf("Dockerfile missing botbox.conf copy:\n%s", dockerfile)
	}
	if !strings.Contains(dockerfile, `CMD ["python3", "src/main.py"]`) {
		t.Errorf("Dockerfile missing plain python CMD:\n%s", dockerfile)
	}
	if strings.Contains(dockerfile, "doppler") {
		t.Errorf("env Dockerfile should not mention doppler:\n%s", dockerfile)
	}

	compose := readOutput(t, filepath.Join(dir, "docker-compose.yml"))
	if !strings.Contains(compose, "env_file:") {
		t.Errorf("env compose missing env_file:\n%s", compose)
	}
	if !strings.Contains(compose, "./logs:/app/logs") {
		t.Errorf("compose missing logs volume:\n%s", compose)
	}

	ignore := readOutput(t, filepath.Join(dir, ".dockerignore"))
	if !strings.Contains(ignore, ".env") || !strings.Contains(ignore, "venv/") {
		t.Errorf(".dockerignore missing expected entries:\n%s", ignore)
	}
}

func TestGenerateDockerFilesDoppler(t *testing.T) {
	dir := t.TempDir()

	if _, err := GenerateDockerFiles(dir, "3.11", "doppler", false); err != nil {
		t.Fatalf("GenerateDockerFiles() error = %v", err)
	}

	dockerfile := readOutput(t, filepath.Join(dir, "Dockerfile"))
	if !strings.Contains(dockerfile, `CMD ["doppler", "run", "--", "python3", "src/main.py"]`) {
		t.Errorf("doppler Dockerfile missing doppler CMD:\n%s", dockerfile)
	}
	if !strings.Contains(dockerfile, "packages.doppler.com") {
		t.Errorf("doppler Dockerfile missing doppler cli install:\n%s", dockerfile)
	}

	compose := readOutput(t, filepath.Join(dir, "docker-compose.yml"))
	if !strings.Contains(compose, "DOPPLER_TOKEN") {
		t.Errorf("doppler compose missing DOPPLER_TOKEN:\n%s", compose)
	}
	if strings.Contains(compose, "env_file:") {
		t.Errorf("doppler compose should not use env_file:\n%s", compose)
	}
}

func TestGenerateDockerFilesSkipsExisting(t *testing.T) {
	dir := t.TempDir()

	// Headless mode makes CreateFileOption skip instead of prompting
	oldHeadless := HeadlessMode
	HeadlessMode = true
	defer func() { HeadlessMode = oldHeadless }()

	if _, err := GenerateDockerFiles(dir, "3.11", "env", false); err != nil {
		t.Fatalf("first GenerateDockerFiles() error = %v", err)
	}

	// A second run without force must leave the existing files untouched
	marker := "# marker"
	dockerfilePath := filepath.Join(dir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(marker), 0644); err != nil {
		t.Fatalf("failed to overwrite Dockerfile: %v", err)
	}

	written, err := GenerateDockerFiles(dir, "3.11", "env", false)
	if err != nil {
		t.Fatalf("second GenerateDockerFiles() error = %v", err)
	}
	if len(written) != 0 {
		t.Errorf("second run wrote %d files, want 0", len(written))
	}
	if got := readOutput(t, dockerfilePath); got != marker {
		t.Errorf("Dockerfile was overwritten without force")
	}

	// Force must overwrite the marker with a rendered Dockerfile again
	written, err = GenerateDockerFiles(dir, "3.11", "env", true)
	if err != nil {
		t.Fatalf("forced GenerateDockerFiles() error = %v", err)
	}
	if len(written) != 3 {
		t.Errorf("forced run wrote %d files, want 3", len(written))
	}
	if got := readOutput(t, dockerfilePath); !strings.Contains(got, "WORKDIR /app") {
		t.Errorf("forced run did not rewrite Dockerfile:\n%s", got)
	}
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
