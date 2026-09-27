package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templates embed.FS

type templateData struct {
	Name     string
	Module   string
	LangID   string
	Features []string
}

func (d templateData) HasFeature(name string) bool {
	return slices.Contains(d.Features, name)
}

func (d templateData) NeedsDocSync() bool {
	return d.HasFeature("hover") || d.HasFeature("diagnostics") ||
		d.HasFeature("completion") || d.HasFeature("definition") ||
		d.HasFeature("references") || d.HasFeature("symbols")
}

var validFeatures = map[string]bool{
	"hover":       true,
	"completion":  true,
	"diagnostics": true,
	"definition":  true,
	"references":  true,
	"formatting":  true,
	"codeactions": true,
	"symbols":     true,
}

var validName = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)

func main() {
	name := flag.String("name", "", "server name")
	module := flag.String("module", "", "Go module path")
	langID := flag.String("lang", "", "language ID")
	features := flag.String("features", "", "comma-separated features (hover,completion,diagnostics,definition,references,formatting,codeactions,symbols)")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	if *name == "" {
		*name = prompt(reader, "Server name")
	}
	if *name == "" {
		fatal("server name is required")
	}
	if !validName.MatchString(*name) {
		fatal("invalid server name %q: use lowercase letters, digits, '-' or '_', starting with a letter", *name)
	}

	if *module == "" {
		*module = prompt(reader, fmt.Sprintf("Module path [github.com/user/%s-lsp]", *name))
		if *module == "" {
			*module = fmt.Sprintf("github.com/user/%s-lsp", *name)
		}
	}

	if *langID == "" {
		*langID = prompt(reader, "Language ID [plaintext]")
		if *langID == "" {
			*langID = "plaintext"
		}
	}

	if *features == "" {
		*features = prompt(reader, "Features (comma-separated) [hover,completion,diagnostics]")
		if *features == "" {
			*features = "hover,completion,diagnostics"
		}
	}

	featureList := parseFeatures(*features)
	for _, f := range featureList {
		if !validFeatures[f] {
			fatal("unknown feature: %s\nvalid features: hover, completion, diagnostics, definition, references, formatting, codeactions, symbols", f)
		}
	}

	data := templateData{
		Name:     *name,
		Module:   *module,
		LangID:   *langID,
		Features: featureList,
	}

	outDir := *name + "-lsp"
	if _, err := os.Stat(outDir); err == nil {
		fatal("%s already exists; remove it or choose another name", outDir)
	}
	if err := generate(outDir, data); err != nil {
		fatal("generation failed: %v", err)
	}

	if _, err := exec.LookPath("go"); err != nil {
		fmt.Fprintln(os.Stderr, "go not found; run 'go get github.com/owenrumney/go-lsp@latest' in the new project")
	} else {
		if err := runGo(outDir, "get", "github.com/owenrumney/go-lsp@latest"); err != nil {
			fmt.Fprintf(os.Stderr, "warning: go get failed: %v\n", err)
		}
		if err := runGo(outDir, "mod", "tidy"); err != nil {
			fmt.Fprintf(os.Stderr, "warning: go mod tidy failed: %v\n", err)
		}
	}

	fmt.Printf("\nCreated %s/\n", outDir)
	fmt.Println("  main.go            — server entrypoint")
	fmt.Println("  handler/handler.go — handler with selected interfaces")
	fmt.Println("  handler/handler_test.go — starter test using servertest")
	fmt.Println("  go.mod")
	fmt.Printf("\nNext: cd %s && go build -o %s .\n", outDir, *name+"-lsp")
}

func generate(outDir string, data templateData) error {
	tmpl, err := template.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		return fmt.Errorf("parsing templates: %w", err)
	}

	files := map[string]string{
		"go.mod":                  "go.mod.tmpl",
		"main.go":                 "main.go.tmpl",
		"handler/handler.go":      "handler.go.tmpl",
		"handler/handler_test.go": "handler_test.go.tmpl",
	}

	for outFile, tmplName := range files {
		path := filepath.Join(outDir, outFile)
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}

		if err := writeTemplate(tmpl, tmplName, path, data); err != nil {
			return fmt.Errorf("writing %s: %w", outFile, err)
		}
	}

	return nil
}

func prompt(reader *bufio.Reader, label string) string {
	fmt.Printf("%s: ", label)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func runGo(dir string, args ...string) error {
	cmd := exec.Command("go", args...) // #nosec G204 -- args are fixed strings
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeTemplate(tmpl *template.Template, name, path string, data templateData) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600) // #nosec G304 -- name is validated
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	return tmpl.ExecuteTemplate(f, name, data)
}

func parseFeatures(s string) []string {
	var features []string
	for f := range strings.SplitSeq(s, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			features = append(features, f)
		}
	}
	return features
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
