package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	root := t.TempDir()
	binary := filepath.Join(root, "import-canonical")
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %s %v", output, err)
	}
	fixtureRoot := filepath.Join("..", "..", "testdata", "historical")
	input := filepath.Join(fixtureRoot, "sections.csv")
	manifest := filepath.Join(fixtureRoot, "manifest.json")
	payload, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	var metadata map[string]any
	if err = json.Unmarshal(payload, &metadata); err != nil {
		t.Fatal(err)
	}
	cutoff := metadata["cutoff"].(string)
	args := []string{"--input", input, "--manifest", manifest, "--cutoff", cutoff}
	run := func(extra []string, want int) string {
		t.Helper()
		output, err := exec.Command(binary, append(append([]string{}, args...), extra...)...).CombinedOutput()
		code := 0
		if err != nil {
			var ok bool
			var exit *exec.ExitError
			exit, ok = err.(*exec.ExitError)
			if !ok {
				t.Fatal(err)
			}
			code = exit.ExitCode()
		}
		if code != want {
			t.Fatalf("exit=%d want=%d: %s", code, want, output)
		}
		return string(output)
	}
	destination := filepath.Join(root, "accepted")
	if output := run([]string{"--output", destination}, 0); !strings.Contains(output, "rows=3 accepted=3 rejected=0") {
		t.Fatal(output)
	}
	run([]string{"--output", destination}, 1)
	run([]string{}, 1)
	run([]string{"--output", filepath.Join(root, "early"), "--cutoff", "2000-01-01T00:00:00Z"}, 1)
	data, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	rejected := filepath.Join(root, "rejected.csv")
	if err = os.WriteFile(rejected, []byte(strings.Replace(string(data), ",100,80,70", ",100,80,81", 1)), 0600); err != nil {
		t.Fatal(err)
	}
	// A modified input cannot use the old provenance hash, even when syntactically valid.
	run([]string{"--input", rejected, "--output", filepath.Join(root, "bad-hash")}, 1)
	// Exercise the quality-specific exit status with an honest new input hash.
	changed, err := os.ReadFile(rejected)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(changed)
	metadata["sha256"] = hex.EncodeToString(digest[:])
	rejectedManifest := filepath.Join(root, "rejected-manifest.json")
	updated, err := json.Marshal(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(rejectedManifest, updated, 0600); err != nil {
		t.Fatal(err)
	}
	rejectedOutput := filepath.Join(root, "quality-rejected")
	run([]string{"--input", rejected, "--manifest", rejectedManifest, "--output", rejectedOutput}, 2)
	if _, err = os.Stat(filepath.Join(rejectedOutput, "normalized.jsonl")); !os.IsNotExist(err) {
		t.Fatal("CLI published a rejected dataset")
	}
	if _, err = os.Stat(filepath.Join(rejectedOutput, "quality.json")); err != nil {
		t.Fatal(err)
	}
}
