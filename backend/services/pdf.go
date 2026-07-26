package services

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

var (
	nonPrintable = regexp.MustCompile(`[^\x20-\x7E\n\t]`)
	multiBlank   = regexp.MustCompile(`\n{3,}`)
	multiSpace   = regexp.MustCompile(` {2,}`)
)

// ExtractTextFromPDF pulls plain text out of an uploaded PDF resume,
// equivalent to the old PyMuPDF-based extract_text_from_pdf().
func ExtractTextFromPDF(fileBytes []byte) (string, error) {
	reader := bytes.NewReader(fileBytes)
	r, err := pdf.NewReader(reader, int64(len(fileBytes)))
	if err != nil {
		return "", fmt.Errorf("failed to parse PDF: %w", err)
	}

	var buf bytes.Buffer
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		buf.WriteString(text)
		buf.WriteString("\n")
	}

	return cleanResumeText(buf.String()), nil
}

func cleanResumeText(text string) string {
	text = nonPrintable.ReplaceAllString(text, " ")
	text = multiBlank.ReplaceAllString(text, "\n\n")
	text = multiSpace.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

// GetResumeSnippet builds the short preview string stored alongside history.
func GetResumeSnippet(text string, maxChars int) string {
	if maxChars <= 0 {
		maxChars = 120
	}
	var lines []string
	for _, l := range strings.Split(text, "\n") {
		l = strings.TrimSpace(l)
		if l != "" {
			lines = append(lines, l)
		}
		if len(lines) == 3 {
			break
		}
	}
	snippet := strings.Join(lines, " · ")
	if len(snippet) > maxChars {
		return snippet[:maxChars] + "…"
	}
	return snippet
}
