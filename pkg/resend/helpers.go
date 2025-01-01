// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package resend

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
)

var templates *template.Template

// LoadTemplates parses embedded mail templates.
func LoadTemplates(templateFS fs.FS) error {
	parsed, err := template.ParseFS(templateFS, "*.html")
	if err != nil {
		return fmt.Errorf("parse mail templates: %w", err)
	}

	templates = parsed
	return nil
}

// renderTemplate renders an email template with data.
func renderTemplate(name string, data any) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return "", fmt.Errorf(
			"render mail template %s: %w",
			name,
			err,
		)
	}

	return buf.String(), nil
}
