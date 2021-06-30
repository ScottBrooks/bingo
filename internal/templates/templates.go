package templates

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var templateFuncs = template.FuncMap{
	"getTime": func() string {
		return time.Now().Format("15:04:05")
	},
	"mod": func(i, j int) bool { return i%j == 0 },
	"isZeroTime": func(t time.Time) bool {
		return t.IsZero()
	},
}

// Template stores the meta data for each template, and whether it uses a layout
type Template struct {
	layout   string
	name     string
	template *template.Template
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	files     embed.FS
	base      string
	templates map[string]*Template
}

// New setup a new template renderer
func New(files embed.FS, base string) *TemplateRenderer {
	return &TemplateRenderer{
		files:     files,
		base:      base,
		templates: make(map[string]*Template),
	}
}

// AddWithLayout register one or more templates using the provided layout
func (t *TemplateRenderer) AddWithLayout(layout string, patterns ...string) error {
	filenames, err := readFileNames(t.files, t.base, patterns...)
	if err != nil {
		return fmt.Errorf("AddWithlayout failed to list using file pattern: %w", err)
	}

	for _, f := range filenames {
		tname := filepath.Base(f)
		lname := filepath.Base(layout)

		log.Debug().Str("filename", tname).Str("layout", layout).Str("tname", tname).Str("lname", lname).Msg("register template")
		t.templates[tname] = &Template{
			layout:   lname,
			name:     tname,
			template: template.Must(template.New(tname).Funcs(templateFuncs).ParseFS(t.files, path.Join(t.base, layout), f)),
		}
	}

	return nil
}

// Add add a template to the registry
func (t *TemplateRenderer) Add(patterns ...string) error {
	filenames, err := readFileNames(t.files, t.base, patterns...)
	if err != nil {
		return fmt.Errorf("failed to read file names using file pattern: %w", err)
	}
	partials, err := readFileNames(t.files, t.base, "partials/*.html")
	if err != nil {
		return fmt.Errorf("failed to read file names using file pattern: %w", err)
	}

	for _, f := range filenames {
		tname := filepath.Base(f)

		log.Debug().Str("filename", f).Msg("register message")
		t.templates[tname] = &Template{
			name:     tname,
			template: template.Must(template.New(tname).Funcs(templateFuncs).ParseFS(t.files, append([]string{f}, partials...)...)),
		}
	}

	return nil
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		log.Ctx(c.Request().Context()).Error().Str("name", name).Msg("template not found")

		return c.NoContent(http.StatusInternalServerError)
	}

	// use the name of the template, or layout if it exists
	execName := tmpl.name
	if tmpl.layout != "" {
		execName = tmpl.layout
	}

	err := tmpl.template.ExecuteTemplate(w, execName, data)
	if err != nil {
		log.Ctx(c.Request().Context()).Error().Err(err).Str("name", tmpl.name).Str("layout", tmpl.layout).Msg("render template failed")
		return err
	}

	return nil
}

func readFileNames(files embed.FS, base string, patterns ...string) ([]string, error) {
	var filenames []string

	for _, pattern := range patterns {
		path := path.Join(base, pattern)
		list, err := fs.Glob(files, path)
		if err != nil {
			return nil, fmt.Errorf("failed to list using file pattern: %w", err)
		}

		if len(list) == 0 {
			return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
		}
		filenames = append(filenames, list...)
	}

	return filenames, nil
}
