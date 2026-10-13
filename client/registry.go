package client

import (
	"encoding/json"
	"slices"
	"strings"
	"sync"

	"github.com/owenrumney/go-lsp/lsp"
)

// registry tracks capabilities the server registered after initialize.
type registry struct {
	mu    sync.RWMutex
	byID  map[string]lsp.Registration
	order []string
}

func newRegistry() *registry {
	return &registry{byID: make(map[string]lsp.Registration)}
}

// Registrations returns the dynamic registrations currently in force for a
// method, in registration order.
func (c *Client) Registrations(method string) []lsp.Registration {
	return c.registry.forMethod(method)
}

// Supports reports whether the server can answer method for uri, considering
// both the capabilities from initialize and anything registered since.
//
// A uri of "" skips document selector matching. Methods outside the set
// covered by this package's typed wrappers are judged on dynamic
// registrations alone, so inspect Capabilities directly for those.
func (c *Client) Supports(method string, uri lsp.DocumentURI, languageID string) bool {
	if c.staticallySupports(method) {
		return true
	}
	for _, reg := range c.registry.forMethod(method) {
		if registrationMatches(reg, uri, languageID) {
			return true
		}
	}
	return false
}

func (r *registry) register(regs []lsp.Registration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, reg := range regs {
		if _, seen := r.byID[reg.ID]; !seen {
			r.order = append(r.order, reg.ID)
		}
		r.byID[reg.ID] = reg
	}
}

func (r *registry) unregister(unregs []lsp.Unregistration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, unreg := range unregs {
		if _, ok := r.byID[unreg.ID]; !ok {
			continue
		}
		delete(r.byID, unreg.ID)
		r.order = slices.DeleteFunc(r.order, func(id string) bool { return id == unreg.ID })
	}
}

func (r *registry) forMethod(method string) []lsp.Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []lsp.Registration
	for _, id := range r.order {
		if reg := r.byID[id]; reg.Method == method {
			out = append(out, reg)
		}
	}
	return out
}

func registrationMatches(reg lsp.Registration, uri lsp.DocumentURI, languageID string) bool {
	if uri == "" || len(reg.RegisterOptions) == 0 {
		return true
	}

	var opts lsp.TextDocumentRegistrationOptions
	if err := json.Unmarshal(reg.RegisterOptions, &opts); err != nil {
		return true
	}
	if opts.DocumentSelector == nil || len(*opts.DocumentSelector) == 0 {
		return true
	}

	for _, filter := range *opts.DocumentSelector {
		if filterMatches(filter, uri, languageID) {
			return true
		}
	}
	return false
}

// filterMatches implements the subset of document filters this package needs.
// Pattern is compared on suffix only; full glob support is not attempted.
func filterMatches(filter lsp.DocumentFilter, uri lsp.DocumentURI, languageID string) bool {
	if filter.Language != "" && filter.Language != languageID {
		return false
	}
	if filter.Scheme != "" && !strings.HasPrefix(string(uri), filter.Scheme+":") {
		return false
	}
	if filter.Pattern != "" {
		ext := strings.TrimPrefix(filter.Pattern, "**/*")
		if strings.HasPrefix(ext, ".") && !suffixMatches(string(uri), ext) {
			return false
		}
	}
	return true
}

// suffixMatches handles a single trailing brace group, as in ".{ts,tsx}".
func suffixMatches(uri, ext string) bool {
	open := strings.IndexByte(ext, '{')
	if open < 0 || !strings.HasSuffix(ext, "}") {
		return strings.HasSuffix(uri, ext)
	}
	prefix := ext[:open]
	for _, alt := range strings.Split(ext[open+1:len(ext)-1], ",") {
		if strings.HasSuffix(uri, prefix+strings.TrimSpace(alt)) {
			return true
		}
	}
	return false
}

func (c *Client) staticallySupports(method string) bool {
	caps := c.caps
	switch method {
	case "textDocument/didOpen", "textDocument/didChange", "textDocument/didClose":
		return caps.TextDocumentSync != nil
	case "textDocument/hover":
		return caps.HoverProvider != nil
	case "textDocument/completion":
		return caps.CompletionProvider != nil
	case "textDocument/definition":
		return caps.DefinitionProvider != nil
	case "textDocument/declaration":
		return caps.DeclarationProvider != nil
	case "textDocument/typeDefinition":
		return caps.TypeDefinitionProvider != nil
	case "textDocument/implementation":
		return caps.ImplementationProvider != nil
	case "textDocument/references":
		return caps.ReferencesProvider != nil
	case "textDocument/documentSymbol":
		// false decodes to non-nil disabled options, so a nil check is not enough.
		return caps.DocumentSymbolProvider.Enabled()
	case "textDocument/codeAction":
		// false decodes to non-nil disabled options, so a nil check is not enough.
		return caps.CodeActionProvider.Enabled()
	case "textDocument/formatting":
		return caps.DocumentFormattingProvider != nil
	case "textDocument/rangeFormatting":
		return caps.DocumentRangeFormattingProvider != nil
	case "textDocument/rename":
		// false decodes to non-nil disabled options, so a nil check is not enough.
		return caps.RenameProvider.Enabled()
	case "textDocument/signatureHelp":
		return caps.SignatureHelpProvider != nil
	case "textDocument/documentHighlight":
		return caps.DocumentHighlightProvider != nil
	case "textDocument/inlayHint":
		// false decodes to non-nil disabled options, so a nil check is not enough.
		return caps.InlayHintProvider.Enabled()
	case "textDocument/diagnostic":
		return caps.DiagnosticProvider != nil
	case "workspace/symbol":
		// false decodes to non-nil disabled options, so a nil check is not enough.
		return caps.WorkspaceSymbolProvider.Enabled()
	case "workspace/executeCommand":
		return caps.ExecuteCommandProvider != nil
	default:
		return false
	}
}
