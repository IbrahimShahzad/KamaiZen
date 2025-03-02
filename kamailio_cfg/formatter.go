package kamailio_cfg

import (
	"KamaiZen/lsp"
	"regexp"
	"strings"
)

func FixIndent(content string) []lsp.TextEdit {
	indentStr := "\t"
	lines := strings.Split(content, "\n")
	var formatted []string
	indentLevel := 0
	edits := []lsp.TextEdit{}
	revertIndent := false

	// Regex to enforce exactly one space before '{'
	braceRegex := regexp.MustCompile(`(\S)\s*\{`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines.
		if trimmed == "" {
			formatted = append(formatted, "")
			continue
		}

		// If the line starts with a closing brace, decrease indent first.
		if strings.HasPrefix(trimmed, "}") {
			indentLevel--
			if indentLevel < 0 {
				// TODO: add diagnostics
				// return "", errors.New("unbalanced braces: too many closing braces")
				return []lsp.TextEdit{}
			}
		}

		// for switch-case statements, decrease indent level
		// if the line starts with 'case' or 'default'
		if strings.HasPrefix(trimmed, "case") || strings.HasPrefix(trimmed, "default") {
			indentLevel--
			revertIndent = true
		}

		// Enforce exactly one space before '{'
		trimmed = braceRegex.ReplaceAllString(trimmed, "$1 {")

		// Build the new indented line.
		currentIndent := strings.Repeat(indentStr, indentLevel)
		formatted = append(formatted, currentIndent+trimmed)

		if revertIndent {
			indentLevel++
			revertIndent = false
		}
		// Increase indent level if line ends with an opening brace.
		if strings.HasSuffix(trimmed, "{") {
			indentLevel++
		}
	}

	if indentLevel != 0 {
		// TODO: add diagnostics
		return []lsp.TextEdit{}
	}
	edit := lsp.TextEdit{
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      0,
				Character: 0,
			},
			End: lsp.Position{
				Line:      len(lines),
				Character: len(lines[len(lines)-1]),
			},
		},
		NewText: strings.Join(formatted, "\n"),
	}
	edits = append(edits, edit)
	return edits
}
