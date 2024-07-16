package routegen

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/andyfusniak/sitebuild/internal/site"
)

type RouteGenerator struct {
	pages map[string]site.Page
}

// NewRouteGenerator creates a new RouteGenerator.
func NewRouteGenerator(pages map[string]site.Page) *RouteGenerator {
	return &RouteGenerator{
		pages: pages,
	}
}

// GenerateRoutes generates the routes for the site.
func (r *RouteGenerator) GenerateRoutes() error {
	routes := `"rewrites": [` + "\n"

	numItems := len(r.pages)
	currentItem := 0
	for key, page := range r.pages {
		currentItem++
		rewrite, err := firebaseRewrite(page.URL, "/"+key)
		if err != nil {
			return err
		}

		if currentItem == numItems {
			routes += rewrite + "\n"
			continue
		}

		routes += rewrite + ",\n"
	}

	routes += `]` + "\n"

	fmt.Println(routes)
	return nil
}

const firebaseRoute = `  {
    "source": "{{.Src}}",
    "destination": "{{.Dst}}"
  }`

var tmpl = template.Must(template.New("firebaseRoute").Parse(firebaseRoute))

func firebaseRewrite(src, dst string) (string, error) {
	tp := struct {
		Src string
		Dst string
	}{
		Src: src,
		Dst: dst,
	}

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, "firebaseRoute", tp); err != nil {
		return "", err
	}
	return buf.String(), nil
}
