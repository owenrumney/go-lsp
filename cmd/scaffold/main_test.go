package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGenerateScaffoldProducesPassingProject(t *testing.T) {
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(t.TempDir(), "mylang-lsp")
	err = generate(outDir, templateData{
		Name:   "mylang",
		Module: "example.com/mylang-lsp",
		LangID: "mylang",
		Features: []string{
			"hover",
			"completion",
			"diagnostics",
			"formatting",
			"codeactions",
			"symbols",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{
		"go.mod",
		"main.go",
		"handler/handler.go",
		"handler/handler_test.go",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected generated %s: %v", name, err)
		}
	}

	modFile := filepath.Join(outDir, "go.mod")
	f, err := os.OpenFile(modFile, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, writeErr := f.WriteString("\nrequire github.com/owenrumney/go-lsp v0.0.0\nreplace github.com/owenrumney/go-lsp => " + filepath.ToSlash(repoRoot) + "\n")
	closeErr := f.Close()
	if writeErr != nil {
		t.Fatal(writeErr)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}

	sum, err := os.ReadFile(filepath.Join(repoRoot, "go.sum"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "go.sum"), sum, 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "mod", "tidy") // #nosec G204 -- fixed command and args
	cmd.Dir = outDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated project tidy failed: %v\n%s", err, output)
	}

	cmd = exec.Command("go", "test", "./...") // #nosec G204 -- fixed command and args
	cmd.Dir = outDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated project tests failed: %v\n%s", err, output)
	}
}

func TestParseFeaturesTrimsEmptyEntries(t *testing.T) {
	got := parseFeatures(" hover, completion,,diagnostics ")
	want := []string{"hover", "completion", "diagnostics"}
	if len(got) != len(want) {
		t.Fatalf("features = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("features = %v, want %v", got, want)
		}
	}
}

func TestValidName(t *testing.T) {
	tests := []struct {
		name string
		ok   bool
	}{
		{"mylang", true},
		{"my-lang_2", true},
		{"", false},
		{"My", false},
		{"../escape", false},
		{"{{.Module}}", false},
		{"a b", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validName.MatchString(tt.name); got != tt.ok {
				t.Fatalf("validName(%q) = %v, want %v", tt.name, got, tt.ok)
			}
		})
	}
}

func TestGenerateRefusesToOverwrite(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "x-lsp")
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "go.mod"), []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := generate(outDir, templateData{Name: "x", Module: "example.com/x", LangID: "x"})
	if err == nil {
		t.Fatal("expected error when file exists")
	}
	got, _ := os.ReadFile(filepath.Join(outDir, "go.mod"))
	if string(got) != "keep" {
		t.Fatalf("existing file overwritten: %q", got)
	}
}
