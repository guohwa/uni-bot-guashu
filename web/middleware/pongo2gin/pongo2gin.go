package pongo2gin

import (
	"errors"
	"net/http"
	"path"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type RenderOptions struct {
	TemplateDir string
	ContentType string
}

type Pongo2Render struct {
	Options  *RenderOptions
	Template *pongo2.Template
	Context  pongo2.Context
}

func init() {
	pongo2.RegisterFilter("unixmilli", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		v, ok := in.Interface().(int64)
		if !ok {
			return nil, &pongo2.Error{
				Sender:    "filter:unixmilli",
				OrigError: errors.New("filter input argument must be of type 'int64'"),
			}
		}
		return pongo2.AsValue(time.UnixMilli(v)), nil
	})
}

func New(options RenderOptions) *Pongo2Render {

	return &Pongo2Render{
		Options: &options,
	}
}

func Default() *Pongo2Render {
	return New(RenderOptions{
		TemplateDir: "templates",
		ContentType: "text/html; charset=utf-8",
	})
}

func (p Pongo2Render) Instance(name string, data interface{}) render.Render {
	var template *pongo2.Template
	filename := path.Join(p.Options.TemplateDir, name)

	if gin.Mode() == "debug" {
		template = pongo2.Must(pongo2.FromFile(filename))
	} else {
		template = pongo2.Must(pongo2.FromCache(filename))
	}

	return Pongo2Render{
		Template: template,
		Context:  data.(pongo2.Context),
		Options:  p.Options,
	}
}

func (p Pongo2Render) Render(w http.ResponseWriter) error {
	p.WriteContentType(w)
	err := p.Template.ExecuteWriter(p.Context, w)
	return err
}

func (p Pongo2Render) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{p.Options.ContentType}
	}
}
