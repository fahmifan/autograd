package web

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Renderer ..
type Renderer struct {
	templates map[string]*template.Template
	location  string
	debug     bool
}

// NewRenderer ..
func NewRenderer(loc string, debug bool) *Renderer {
	rd := &Renderer{location: loc, templates: make(map[string]*template.Template), debug: debug}
	rd.parseTemplates()
	return rd
}

func (r *Renderer) parse(path string) error {
	tpl, err := template.ParseFiles(path)
	if err != nil {
		log.Println(err)
		return err
	}

	// generally we dont need to locks
	// because in production we only did this once
	r.templates[path] = tpl

	return nil
}

func (r *Renderer) parseTemplates() error {
	err := filepath.Walk(r.location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Error(err)
			return err
		}

		if !strings.Contains(path, ".html") {
			return nil
		}

		if err = r.parse(path); err != nil {
			log.Println(err)
		}
		return err
	})

	return err
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	realPath := path.Join(r.location, name)
	if r.debug {
		if err := r.parse(realPath); err != nil {
			logrus.Error(err)
		}
	}

	tpl, ok := r.templates[realPath]
	if !ok {
		logrus.Error("unable to find " + realPath)
		return echo.ErrNotFound
	}

	return tpl.Execute(w, data)
}
